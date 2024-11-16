// simulation.go
package main

func (sim *Simulation) ProcessNextStep() string {
	if sim.CurrentStep >= len(sim.PageSequence) {
		sim.Status = "Completed"
		return "Simulation completed."
	}

	page := sim.PageSequence[sim.CurrentStep]
	sim.CurrentStep++

	// Check if page is already in memory
	for _, p := range sim.MemoryFrames {
		if p == page {
			sim.PageHits++
			if sim.Algorithm == "LRU" {
				sim.LRUTracker[page] = sim.CurrentStep
			}
			return "Page hit."
		}
	}

	// Page fault occurs
	sim.PageFaults++

	// If there is empty frame
	if len(sim.MemoryFrames) < sim.NumFrames {
		sim.MemoryFrames = append(sim.MemoryFrames, page)
	} else {
		// Replace a page based on algorithm
		var replaced int
		if sim.Algorithm == "FIFO" {
			replaced = sim.FIFOQueue[0]
			sim.FIFOQueue = sim.FIFOQueue[1:]
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

	return "Page fault."
}
