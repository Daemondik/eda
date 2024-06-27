package main

import (
	"context"
	"crypto/rand"
	"eda/internal/sms"
	"eda/logger"
	"eda/models"
	"eda/utils/security"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"math/big"
	"net"
	"os"
	"os/signal"
	"strconv"
	"time"

	pb "eda/internal/auth"
	"google.golang.org/grpc"
)

type Config struct {
	CertFile string
	KeyFile  string
}

type server struct {
	pb.UnimplementedAuthServiceServer
	config *Config
}

func NewConfig() *Config {
	return &Config{
		CertFile: os.Getenv("CERT_FILE"),
		KeyFile:  os.Getenv("KEY_FILE"),
	}
}

func (s *server) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	// Создаем экземпляр LoginInput и заполняем его данными из запроса
	input := models.LoginInput{
		Phone:    in.GetPhone(),
		Password: in.GetPassword(),
	}

	// Проверяем валидность JSON
	if err := validateLoginInput(input); err != nil {
		return nil, err
	}

	// Проверяем валидность номера телефона
	if phoneValid := security.IsValidRussianPhoneNumber(input.Phone); !phoneValid {
		return nil, fmt.Errorf("phone should be format 7XXXXXXXXXX")
	}

	// Проверяем учетные данные и генерируем токен
	generatedToken, err := models.LoginCheck(input.Phone, input.Password)
	if err != nil {
		return nil, fmt.Errorf("phone or password is incorrect")
	}

	// Возвращаем успешный ответ с токеном
	return &pb.LoginResponse{
		Token: generatedToken,
	}, nil
}

// Register - gRPC метод для регистрации пользователя.
func (s *server) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	u := models.User{
		Phone:    in.Phone,
		Password: in.Password,
		IsActive: false,
	}

	// Проверка номера телефона.
	if phoneValid := security.IsValidRussianPhoneNumber(u.Phone); !phoneValid {
		return nil, status.Errorf(codes.InvalidArgument, "phone should be format 7XXXXXXXXXX")
	}

	// Генерация кода подтверждения.
	n, err := rand.Int(rand.Reader, big.NewInt(9000))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "code generation error: %s", err.Error())
	}
	code := n.Int64() + 1000

	// Отправка SMS кода
	smsSender := sms.NewSMSruSender()
	err = smsSender.SendSMSCode(u.Phone, strconv.FormatInt(code, 10))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "sending error: %s", err.Error())
	}

	// Сохранение пользователя.
	_, err = u.SaveUser()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error saving user: %s", err.Error())
	}

	// Установка кода в Redis.
	expiration := time.Now().Add(time.Hour)
	models.RedisClient.Set(u.Phone, code, expiration.Sub(time.Now()))

	return &pb.RegisterResponse{}, nil
}

// ConfirmSMSCode - gRPC метод для подтверждения SMS-кода.
func (s *server) ConfirmSMSCode(ctx context.Context, in *pb.ConfirmSMSCodeRequest) (*pb.ConfirmSMSCodeResponse, error) {
	// Получение текущего кода из Redis.
	currentCode, err := models.GetDelPhoneTransaction(in.Phone)
	if err != nil {
		logger.Log.Error("get code: " + err.Error() + "\n")
		return nil, status.Errorf(codes.Internal, "get code: %s", err.Error())
	}

	// Проверка кода.
	if in.Code != currentCode {
		return nil, status.Errorf(codes.InvalidArgument, "incorrect code")
	}

	// Получение пользователя по номеру телефона.
	u, err := models.GetUserByPhone(in.Phone)
	if err != nil {
		logger.Log.Error("User Exist: " + err.Error() + "\n")
		return nil, status.Errorf(codes.NotFound, "user not found: %s", err.Error())
	}

	// Активация пользователя.
	_, err = u.SetActive()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error setting user active: %s", err.Error())
	}

	return &pb.ConfirmSMSCodeResponse{Message: "registration success"}, nil
}

func main() {
	config := NewConfig()

	if err := models.InitializeServices(); err != nil {
		logger.Log.Fatal("Failed to initialize services: ", zap.Error(err))
	}

	creds, err := credentials.NewServerTLSFromFile(config.CertFile, config.KeyFile)
	if err != nil {
		logger.Log.Fatal("Failed to generate credentials", zap.Error(err))
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Log.Fatal("failed to listen", zap.Error(err))
	}

	s := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterAuthServiceServer(s, &server{config: config})

	go func() {
		if err := s.Serve(lis); err != nil {
			logger.Log.Fatal("failed to serve", zap.Error(err))
		}
	}()

	waitForShutdown(s)
}

// Ожидание сигнала для graceful shutdown.
func waitForShutdown(srv *grpc.Server) {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)
	<-stopChan

	srv.GracefulStop()
	logger.Log.Info("shutting down server")
}

func validateLoginInput(input models.LoginInput) error {
	// TODO: сделать дополнительную логику валидации
	if input.Phone == "" || input.Password == "" {
		return fmt.Errorf("phone and password are required")
	}
	return nil
}
