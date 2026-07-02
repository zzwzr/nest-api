package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func RandomPassword(length int) (string, error) {
	if length <= 0 {
		length = 24
	}

	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes)[:length], nil
}
