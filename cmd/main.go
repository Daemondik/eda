package main

import (
	"eda/logger"
	"eda/middlewares"
	"eda/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
)

func main() {
	err := logger.InitializeZapCustomLogger()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer func(Log *zap.Logger) {
		err := Log.Sync()
		if err != nil {
			log.Fatal("Failed to sync logger:", err)
		}
	}(logger.Log)

	models.ConnectDb()

	//gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	public := r.Group("/api")
	setupApiRoutes(public)

	protected := r.Group("/api/admin")
	protected.Use(middlewares.AuthMiddleware(models.RoleAdmin))
	setupApiAdminRoutes(protected)

	err = r.Run()
	if err != nil {
		return
	}
}
