package handlers

import (
	"fmt"
	"net/http"
	"memory-simulation/models"
	"memory-simulation/utils"

	"github.com/gin-gonic/gin"
)

func Simulate(c *gin.Context) {
	var req models.SimulationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var pr models.PageReplacement
	switch req.Algorithm {
	case "FIFO":
		pr = utils.NewFIFO(req.NumFrames)
	case "MRU":
		pr = utils.NewMRU(req.NumFrames)
	case "LRU":
		pr = utils.NewLRU(req.NumFrames)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported algorithm"})
		return
	}

	response := models.SimulationResponse{}
	physicalFrames := make(map[int]int) // Page -> Frame
	frameToPage := make(map[int]int)     // Frame -> Page
	nextFrame := 0

	for i, page := range req.Pages {
		hit, action := pr.Access(page)
		step := models.SimulationStep{
			Step: i + 1,
			Page: page,
			Hit:  hit,
		}

		if hit {
			step.Action = "Page Hit"
		} else {
			if action[:7] == "Evicted" {
				step.Action = action + "; Loaded into frame " + fmt.Sprint(nextFrame)
				// Find which frame had the evicted page
				// For simplicity, assume frame allocation follows nextFrame
				// In a real implementation, track frame allocations more accurately
				nextFrame = (nextFrame + 1) % req.NumFrames
			} else {
				step.Action = action + " " + fmt.Sprint(nextFrame)
				nextFrame = (nextFrame + 1) % req.NumFrames
			}
			physicalFrames[page] = nextFrame
			frameToPage[nextFrame] = page
		}

		// Translate virtual address to physical address
		virtualAddress := req.VirtualAddresses[i]
		frameNumber, exists := physicalFrames[page]
		if !exists {
			frameNumber = nextFrame
		}
		physicalAddress := fmt.Sprintf("0x%X", (frameNumber<<12)|virtualAddress.Offset)
		step.PhysicalAddress = physicalAddress

		response.SimulationSteps = append(response.SimulationSteps, step)
		if hit {
			response.TotalHits++
		} else {
			response.TotalMisses++
		}
	}

	c.JSON(http.StatusOK, response)
}
