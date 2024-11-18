package models

type VirtualAddress struct {
	Page   int `json:"page"`
	Offset int `json:"offset"`
}

type SimulationRequest struct {
	NumFrames       int             `json:"numFrames" binding:"required"`
	Pages           []int           `json:"pages" binding:"required"`
	VirtualAddresses []VirtualAddress `json:"virtualAddresses" binding:"required"`
	Algorithm       string          `json:"algorithm" binding:"required,oneof=FIFO MRU LRU"`
}

type SimulationStep struct {
	Step           int    `json:"step"`
	Page           int    `json:"page"`
	Hit            bool   `json:"hit"`
	Action         string `json:"action"`
	PhysicalAddress string `json:"physicalAddress"`
}

type SimulationResponse struct {
	SimulationSteps []SimulationStep `json:"simulationSteps"`
	TotalHits       int               `json:"totalHits"`
	TotalMisses     int               `json:"totalMisses"`
}
