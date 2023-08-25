package generator

import (
	"encoding/base64"
	"math/rand"
)

type Default struct {
}

func (s *Default) Rand() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 20

	b := make([]byte, keyLength)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return base64.StdEncoding.EncodeToString(b)
}
