package main

import (
	"eda/handlers"
	"github.com/gin-gonic/gin"
)

func setupApiRoutes(app *gin.RouterGroup) {
	app.GET("/fact-list", handlers.ListFacts)
	app.POST("/create-fact", handlers.CreateFact)

	app.POST("/register", handlers.Register)
	app.POST("/confirm-sms", handlers.ConfirmSMSCode)
	app.POST("/login", handlers.Login)

	app.POST("/login-gl", handlers.GoogleLogin)
	app.GET("/callback-gl", handlers.CallBackFromGoogle)
}

func setupApiAdminRoutes(app *gin.RouterGroup) {
	app.GET("/user", handlers.CurrentUser)
	//app.GET("/profile", handlers.Profile)
}
