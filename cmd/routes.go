package main

import (
	"eda/handlers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func setupApiRoutes(app *gin.RouterGroup) {
	app.POST("/register", handlers.Register)
	app.POST("/confirm-sms", handlers.ConfirmSMSCode)
	app.POST("/login", handlers.Login)

	app.POST("/login-gl", handlers.GoogleLogin)
	app.GET("/callback-gl", handlers.CallBackFromGoogle)
}

func setupApiAdminRoutes(app *gin.RouterGroup) {
	app.GET("/user", handlers.CurrentUser)
}

func setupWebsocketRoutes(app *gin.RouterGroup) {
	app.GET("/chat/:recipient_id", handlers.Chat)
}

func setupFrontRoutes(app *gin.RouterGroup) {
	app.GET("/chat/:user_id", func(c *gin.Context) {
		c.HTML(http.StatusOK, "chat.tmpl", gin.H{
			"user_id": c.Param("user_id"),
		})
	})
}
