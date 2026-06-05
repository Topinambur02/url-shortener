package utils

import "crypto/sha256"

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	length   = 10
)

func GenerateShortUrl(originalUrl string) string {
	hash := sha256.Sum256([]byte(originalUrl))
	var result []byte
	base := len(alphabet)

	for i := range length {
		idx := int(hash[i]) % base
		result = append(result, alphabet[idx])
	}

	return string(result)
}
