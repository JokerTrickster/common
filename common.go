package common

import (
	"time"

	"math/rand"
)

// GenerateRandomFilename generates a random filename of the given length
func GenerateRandomFilename(length int) string {
	// Validate input length
	if length <= 0 {
		return ""
	}

	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rand.Seed(time.Now().UnixNano())

	// Create a slice of the specified length
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
