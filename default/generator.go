package _default

import (
	"encoding/base64"
	"math/rand"
)

type Generator struct {
}

func (s *Generator) Rand() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 20

	b := make([]byte, keyLength)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return base64.StdEncoding.EncodeToString(b)
}
