package handlers

import (
	"context"
	pb "eda/internal/auth"
	"eda/logger"
	"eda/models"
	"eda/utils/security"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net/http"
	"os"
)

func CurrentUser(c *gin.Context) {
	userId, err := security.GetUserIdByJWTOrOauth(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := models.GetUserByID(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": u})
}

func createGRPCClient() (pb.AuthServiceClient, func(), error) {
	certFile := os.Getenv("CERT_FILE")
	creds, err := credentials.NewClientTLSFromFile(certFile, "")
	if err != nil {
		return nil, nil, err
	}

	conn, err := grpc.Dial("host.docker.internal:50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, nil, err
	}

	client := pb.NewAuthServiceClient(conn)
	return client, func() { conn.Close() }, nil
}

func Login(c *gin.Context) {
	var input models.LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	client, closeConn, err := createGRPCClient()
	if err != nil {
		logger.Log.Error("Failed to create gRPC client", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to gRPC server"})
		return
	}
	defer closeConn()

	response, err := client.Login(c, &pb.LoginRequest{
		Phone:    input.Phone,
		Password: input.Password,
	})
	if err != nil {
		logger.Log.Error("Authentication failed", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": response.Token})
}

func Register(c *gin.Context) {
	var input models.RegisterInput
	var err error

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	client, closeConn, err := createGRPCClient()
	if err != nil {
		logger.Log.Error("Failed to create gRPC client", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to gRPC server"})
		return
	}
	defer closeConn()

	_, err = client.Register(context.Background(), &pb.RegisterRequest{
		Phone:    input.Phone,
		Password: input.Password,
	})
	if err != nil {
		logger.Log.Error("Register failed", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Register failed"})
		return
	}
}

func ConfirmSMSCode(c *gin.Context) {
	var input models.ConfirmSMSCodeInput
	var err error

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, closeConn, err := createGRPCClient()
	if err != nil {
		logger.Log.Error("Failed to create gRPC client", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to gRPC server"})
		return
	}
	defer closeConn()

	_, err = client.ConfirmSMSCode(context.Background(), &pb.ConfirmSMSCodeRequest{
		Phone: input.Phone,
		Code:  input.Code,
	})
	if err != nil {
		logger.Log.Error("Register failed", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Register failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration success"})
}
