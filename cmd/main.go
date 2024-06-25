package main

import (
	"eda/logger"
	"eda/middlewares"
	"eda/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	if err := models.InitializeServices(); err != nil {
		logger.Log.Fatal("Failed to initialize services: ", zap.Error(err))
	}

	r := setupRouter()
	if err := r.Run(); err != nil {
		logger.Log.Fatal("Failed to run the server:", zap.Error(err))
	}
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	public := r.Group("/api")
	setupApiRoutes(public)

	protected := r.Group("/api/admin")
	protected.Use(middlewares.AuthMiddleware(models.Admin))
	setupApiAdminRoutes(protected)

	ws := r.Group("/ws")
	setupWebsocketRoutes(ws)

	r.LoadHTMLGlob("front/templates/*")
	front := r.Group("/")
	setupFrontRoutes(front)

	return r
}
