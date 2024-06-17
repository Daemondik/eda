package main

import (
	"eda/logger"
	"eda/middlewares"
	"eda/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"os"
)

func main() {
	err := logger.InitializeZapCustomLogger()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	if os.Getenv("LOGGER_OUTPUT_PATH") != "" {
		defer func(Log *zap.Logger) {
			err := Log.Sync()
			if err != nil {
				log.Printf("Failed to sync logger:", err)
			}
		}(logger.Log)
	}

	models.ConnectDb()
	models.NewRedis()

	//gin.SetMode(gin.ReleaseMode)
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

	err = r.Run()
	if err != nil {
		return
	}
}
