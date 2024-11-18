package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

// GetAlgorithms returns the list of available page replacement algorithms.
func GetAlgorithms(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "algorithms": []string{"FIFO", "MRU", "LRU"},
    })
}
