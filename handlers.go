// handlers.go
package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type InitRequest struct {
	NumFrames    int   `json:"num_frames" binding:"required"`
	Algorithm    string `json:"algorithm" binding:"required,oneof=FIFO LRU"`
	PageSequence []int  `json:"page_sequence" binding:"required"`
}

func initSimulation(c *gin.Context) {
	var req InitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := generateID()

	sim := &Simulation{
		ID:           id,
		NumFrames:    req.NumFrames,
		Algorithm:    req.Algorithm,
		PageSequence: req.PageSequence,
		MemoryFrames: []int{},
		PageFaults:   0,
		PageHits:     0,
		CurrentStep:  0,
		StartTime:    time.Now(),
		Status:       "In Progress",
		LRUTracker:   make(map[int]int),
		FIFOQueue:    []int{},
	}

	simLock.Lock()
	simulations[id] = sim
	simLock.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"message":        "Memory simulation initialized successfully.",
		"simulation_id": id,
	})
}

func getSimulationState(c *gin.Context) {
	id := c.Param("id")
	simLock.RLock()
	sim, exists := simulations[id]
	simLock.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Simulation not found."})
		return
	}

	sim.Lock.RLock()
	defer sim.Lock.RUnlock()

	c.JSON(http.StatusOK, gin.H{
		"current_step":  sim.CurrentStep,
		"memory_frames": sim.MemoryFrames,
		"page_faults":   sim.PageFaults,
		"page_hits":     sim.PageHits,
		"status":        sim.Status,
	})
}

func advanceSimulation(c *gin.Context) {
	id := c.Param("id")
	simLock.RLock()
	sim, exists := simulations[id]
	simLock.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Simulation not found."})
		return
	}

	sim.Lock.Lock()
	defer sim.Lock.Unlock()

	if sim.Status == "Completed" {
		c.JSON(http.StatusOK, gin.H{
			"message": "Simulation already completed.",
		})
		return
	}

	action := sim.ProcessNextStep()

	c.JSON(http.StatusOK, gin.H{
		"current_step":  sim.CurrentStep,
		"memory_frames": sim.MemoryFrames,
		"page_faults":   sim.PageFaults,
		"page_hits":     sim.PageHits,
		"action":        action,
		"status":        sim.Status,
	})
}

func resetSimulation(c *gin.Context) {
	id := c.Param("id")
	simLock.RLock()
	sim, exists := simulations[id]
	simLock.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Simulation not found."})
		return
	}

	sim.Lock.Lock()
	defer sim.Lock.Unlock()

	sim.MemoryFrames = []int{}
	sim.PageFaults = 0
	sim.PageHits = 0
	sim.CurrentStep = 0
	sim.Status = "In Progress"
	sim.LRUTracker = make(map[int]int)
	sim.FIFOQueue = []int{}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Simulation reset successfully.",
		"current_step":  sim.CurrentStep,
		"memory_frames": sim.MemoryFrames,
		"page_faults":   sim.PageFaults,
		"page_hits":     sim.PageHits,
		"status":        sim.Status,
	})
}

func getSimulationResults(c *gin.Context) {
	id := c.Param("id")
	simLock.RLock()
	sim, exists := simulations[id]
	simLock.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Simulation not found."})
		return
	}

	sim.Lock.RLock()
	defer sim.Lock.RUnlock()

	c.JSON(http.StatusOK, gin.H{
		"total_steps":        sim.CurrentStep,
		"total_page_faults":  sim.PageFaults,
		"total_page_hits":    sim.PageHits,
		"algorithm":          sim.Algorithm,
		"page_sequence":      sim.PageSequence,
		"memory_frames_final": sim.MemoryFrames,
	})
}

func listSimulations(c *gin.Context) {
	simLock.RLock()
	defer simLock.RUnlock()

	var list []gin.H
	for _, sim := range simulations {
		list = append(list, gin.H{
			"simulation_id": sim.ID,
			"algorithm":     sim.Algorithm,
			"status":        sim.Status,
		})
	}

	c.JSON(http.StatusOK, list)
}
