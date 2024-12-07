package utils

import (
    "math/rand"
    "time"
)

// GenerateFamilyToken generates a unique 8-character token
func GenerateFamilyToken() string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    const length = 8

    rand.Seed(time.Now().UnixNano()) // Ensure randomness
    token := make([]byte, length)
    for i := range token {
        token[i] = charset[rand.Intn(len(charset))]
    }
    return string(token)
}