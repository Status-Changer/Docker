package utils

import (
	"math/rand"
	"time"
)

func RandStringBytes(n int) string {
	letterBytes := "1234567890abcdef"
	rand.Seed(time.Now().UnixNano())
	randomBytes := make([]byte, n)
	for i := range randomBytes {
		randomBytes[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(randomBytes)
}
