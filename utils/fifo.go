package utils

import (
	"container/list"
	"fmt"
)

type FIFO struct {
	capacity int
	cache    map[int]bool
	queue    *list.List
}

// GetPhysicalAddress implements models.PageReplacement.
func (f *FIFO) GetPhysicalAddress(page int, offset int) string {
	panic("unimplemented")
}

func NewFIFO(capacity int) *FIFO {
	return &FIFO{
		capacity: capacity,
		cache:    make(map[int]bool),
		queue:    list.New(),
	}
}

func (f *FIFO) Access(page int) (hit bool, action string) {
	if _, exists := f.cache[page]; exists {
		return true, "Page Hit"
	}
	// Page Miss
	if f.queue.Len() >= f.capacity {
		// Evict the oldest page
		evicted := f.queue.Remove(f.queue.Front()).(int)
		delete(f.cache, evicted)
		action = "Evicted page " + fmt.Sprint(evicted)
	}
	f.queue.PushBack(page)
	f.cache[page] = true
	if action == "" {
		action = "Loaded into frame"
	}
	return false, action
}
