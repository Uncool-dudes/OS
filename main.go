// main.go
package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.POST("/init", initSimulation)
		api.GET("/simulations", listSimulations)

		sim := api.Group("/simulation/:id")
		{
			sim.GET("/state", getSimulationState)
			sim.POST("/next", advanceSimulation)
			sim.POST("/reset", resetSimulation)
			sim.GET("/results", getSimulationResults)
		}
	}

	// Serve frontend (if using Go templates)
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	router.Run(":8080")
}
