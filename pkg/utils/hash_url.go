package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashURL(inputURL string) string {
	hasher := sha256.New()
	hasher.Write([]byte(inputURL))
	hashBytes := hasher.Sum(nil)

	return hex.EncodeToString(hashBytes)
}