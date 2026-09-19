package main

import "math/rand/v2"

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateCode(n int) string {
	result := make([]byte, n)
	for i := 0; i < n; i++ {
		idx := rand.IntN(len(alphabet))
		result[i] = alphabet[idx]
	}
	return string(result)
}
