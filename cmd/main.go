package main

import (
	"eda/middlewares"
	"eda/models"
	"github.com/gin-gonic/gin"
)

func main() {
	models.ConnectDb()

	//gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	public := r.Group("/api")
	setupApiRoutes(public)

	protected := r.Group("/api/admin")
	protected.Use(middlewares.JwtAuthMiddleware(models.RoleAdmin))
	setupApiAdminRoutes(protected)

	err := r.Run()
	if err != nil {
		return
	}
}
