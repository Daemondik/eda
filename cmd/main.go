package main

import (
	"eda/database"
	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectDb()

	r := gin.Default()
	public := r.Group("/api")

	setupApiRoutes(public)

	err := r.Run()
	if err != nil {
		return
	}
}
