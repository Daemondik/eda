package main

import (
	"eda/logger"
	"eda/middlewares"
	"eda/models"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	err := logger.InitializeZapCustomLogger()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
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

	r.LoadHTMLGlob("../front/templates/*")
	front := r.Group("/")
	setupFrontRoutes(front)

	err = r.Run()
	if err != nil {
		return
	}
}
