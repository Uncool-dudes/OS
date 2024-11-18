package main

import (
	"memory-simulation/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Serve static files
	router.Static("/static", "./static")

	// Serve HTML templates
	router.LoadHTMLGlob("templates/*")

	// API Routes
	api := router.Group("/api")
	{
		api.GET("/algorithms", handlers.GetAlgorithms)
		api.POST("/simulate", handlers.Simulate)
	}

	// Frontend Routes
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	router.Run(":8080")
}
