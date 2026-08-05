package utils

import "math/rand/v2"

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateShortCode() string {
	code := make([]byte, 6)
	for i := range code {
		idx := rand.IntN(len(charset))
		code[i] = charset[idx]
	}
	return string(code)
}
