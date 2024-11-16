// simulation.go
package main

import "fmt"

func (sim *Simulation) ProcessNextStep() string {
	if sim.CurrentStep >= len(sim.PageSequence) {
		sim.Status = "Completed"
		return "Simulation completed."
	}

	page := sim.PageSequence[sim.CurrentStep]
	sim.CurrentStep++

	// Initialize action description
	actionDesc := ""

	// Check if page is already in memory
	for _, p := range sim.MemoryFrames {
		if p == page {
			sim.PageHits++
			actionDesc = "Page hit."
			if sim.Algorithm == "LRU" {
				sim.LRUTracker[page] = sim.CurrentStep
			}
			// Record the step
			sim.Steps = append(sim.Steps, Step{
				StepNumber:     sim.CurrentStep,
				VirtualAddress: page,
				PhysicalFrames: append([]int{}, sim.MemoryFrames...), // Make a copy
				Action:         actionDesc,
			})
			return actionDesc
		}
	}

	// Page fault occurs
	sim.PageFaults++
	actionDesc = "Page fault."

	// If there is empty frame
	if len(sim.MemoryFrames) < sim.NumFrames {
		sim.MemoryFrames = append(sim.MemoryFrames, page)
		actionDesc += " Loaded into empty frame."
	} else {
		// Replace a page based on algorithm
		var replaced int
		if sim.Algorithm == "FIFO" {
			replaced = sim.FIFOQueue[0]
			sim.FIFOQueue = sim.FIFOQueue[1:]
			actionDesc += " Replaced page " + itoa(replaced) + " using FIFO."
		} else if sim.Algorithm == "LRU" {
			// Find the least recently used page
			minTime := sim.CurrentStep
			var lruPage int
			for p, t := range sim.LRUTracker {
				if t < minTime {
					minTime = t
					lruPage = p
				}
			}
			replaced = lruPage
			delete(sim.LRUTracker, replaced)
			actionDesc += " Replaced page " + itoa(replaced) + " using LRU."
		}

		// Replace the page
		for i, p := range sim.MemoryFrames {
			if p == replaced {
				sim.MemoryFrames[i] = page
				break
			}
		}

		if sim.Algorithm == "FIFO" {
			sim.FIFOQueue = append(sim.FIFOQueue, page)
		} else if sim.Algorithm == "LRU" {
			sim.LRUTracker[page] = sim.CurrentStep
		}
	}

	// Update FIFO queue
	if sim.Algorithm == "FIFO" && sim.CurrentStep <= len(sim.PageSequence) {
		sim.FIFOQueue = append(sim.FIFOQueue, page)
	}

	// Record the step
	sim.Steps = append(sim.Steps, Step{
		StepNumber:     sim.CurrentStep,
		VirtualAddress: page,
		PhysicalFrames: append([]int{}, sim.MemoryFrames...), // Make a copy
		Action:         actionDesc,
	})

	return actionDesc
}

// Helper function to convert int to string
func itoa(num int) string {
	return fmt.Sprintf("%d", num)
}
