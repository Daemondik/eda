package main

import (
	"eda/logger"
	"eda/middlewares"
	"eda/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	if err := initializeServices(); err != nil {
		logger.Log.Fatal("Failed to initialize services: ", zap.Error(err))
	}

	r := setupRouter()
	if err := r.Run(); err != nil {
		logger.Log.Fatal("Failed to run the server:", zap.Error(err))
	}
}

func initializeServices() error {
	if err := logger.InitializeZapCustomLogger(); err != nil {
		return err
	}

	if err := models.ConnectDb(); err != nil {
		return err
	}

	if err := models.NewRedis(); err != nil {
		return err
	}

	return nil
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	public := r.Group("/api")
	setupApiRoutes(public)

	protected := r.Group("/api/admin")
	protected.Use(middlewares.AuthMiddleware(models.RoleAdmin))
	setupApiAdminRoutes(protected)

	ws := r.Group("/ws")
	setupWebsocketRoutes(ws)

	r.LoadHTMLGlob("front/templates/*")
	front := r.Group("/")
	setupFrontRoutes(front)

	return r
}
