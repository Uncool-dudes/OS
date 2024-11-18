package utils

import (
	"container/list"
	"fmt"
)

type MRU struct {
	capacity int
	cache    map[int]*list.Element
	order    *list.List
}

// GetPhysicalAddress implements models.PageReplacement.
func (m *MRU) GetPhysicalAddress(page int, offset int) string {
	panic("unimplemented")
}

type mruEntry struct {
	page int
}

func NewMRU(capacity int) *MRU {
	return &MRU{
		capacity: capacity,
		cache:    make(map[int]*list.Element),
		order:    list.New(),
	}
}

func (m *MRU) Access(page int) (hit bool, action string) {
	if elem, ok := m.cache[page]; ok {
		m.order.MoveToFront(elem)
		return true, "Page Hit"
	}
	// Page Miss
	if m.order.Len() >= m.capacity {
		// Evict the most recently used
		evicted := m.order.Front()
		if evicted != nil {
			m.order.Remove(evicted)
			delete(m.cache, evicted.Value.(mruEntry).page)
			action = fmt.Sprintf("Evicted page %d", evicted.Value.(mruEntry).page)
		}
	}
	elem := m.order.PushFront(mruEntry{page: page})
	m.cache[page] = elem
	if action == "" {
		action = "Loaded into frame"
	}
	return false, action
}
