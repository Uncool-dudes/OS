package models

type PageReplacement interface {
	Access(page int) (hit bool, action string)
	GetPhysicalAddress(page int, offset int) string
}
