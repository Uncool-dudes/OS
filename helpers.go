// helpers.go
package main

import (
	"math/rand"
	"time"
)

func generateID() string {
	rand.Seed(time.Now().UnixNano())
	return "sim-" + randSeq(8)
}

func randSeq(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
