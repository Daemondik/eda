package main

import (
	"eda/handlers"
	"github.com/gin-gonic/gin"
)

func setupApiRoutes(app *gin.RouterGroup) {
	app.GET("/fact-list", handlers.ListFacts)
	app.POST("/create-fact", handlers.CreateFact)

	app.POST("/register", handlers.Register)
}
