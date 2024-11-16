// models.go
package main

import (
	"sync"
	"time"
)

type Step struct {
	StepNumber     int   `json:"step_number"`
	VirtualAddress int   `json:"virtual_address"`
	PhysicalFrames []int `json:"physical_frames"`
	Action         string `json:"action"`
}

type Simulation struct {
	ID           string        `json:"simulation_id"`
	NumFrames    int           `json:"num_frames"`
	Algorithm    string        `json:"algorithm"` // "FIFO" or "LRU"
	PageSequence []int         `json:"page_sequence"`
	MemoryFrames []int         `json:"memory_frames"`
	PageFaults   int           `json:"page_faults"`
	PageHits     int           `json:"page_hits"`
	CurrentStep  int           `json:"current_step"`
	StartTime    time.Time     `json:"start_time"`
	Status       string        `json:"status"` // "In Progress", "Completed"
	Steps        []Step        `json:"steps"`
	Lock         sync.RWMutex  `json:"-"`
	LRUTracker   map[int]int   `json:"-"`
	FIFOQueue    []int         `json:"-"`
}

var simulations = make(map[string]*Simulation)
var simLock = sync.RWMutex{}
