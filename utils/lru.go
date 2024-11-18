package utils

import (
	"container/list"
	"fmt"
)

type LRU struct {
	capacity int
	cache    map[int]*list.Element
	order    *list.List
}

// GetPhysicalAddress implements models.PageReplacement.
func (l *LRU) GetPhysicalAddress(page int, offset int) string {
	panic("unimplemented")
}

type lruEntry struct {
	page int
}

func NewLRU(capacity int) *LRU {
	return &LRU{
		capacity: capacity,
		cache:    make(map[int]*list.Element),
		order:    list.New(),
	}
}

func (l *LRU) Access(page int) (hit bool, action string) {
	if elem, ok := l.cache[page]; ok {
		l.order.MoveToFront(elem)
		return true, "Page Hit"
	}
	// Page Miss
	if l.order.Len() >= l.capacity {
		// Evict least recently used
		evicted := l.order.Back()
		if evicted != nil {
			l.order.Remove(evicted)
			delete(l.cache, evicted.Value.(lruEntry).page)
			action = fmt.Sprintf("Evicted page %d", evicted.Value.(lruEntry).page)
		}
	}
	elem := l.order.PushFront(lruEntry{page: page})
	l.cache[page] = elem
	if action == "" {
		action = "Loaded into frame"
	}
	return false, action
}
