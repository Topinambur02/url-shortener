package utils

import "hash/fnv"

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	length   = 10
)

func GenerateShortUrl(originalUrl string) string {
	h := fnv.New64a()
	h.Write([]byte(originalUrl))
	hashValue := h.Sum64()
	var result [length]byte
	base := uint64(len(alphabet))

	for i := range length {
		result[i] = alphabet[hashValue%base]
		hashValue /= base
	}

	return string(result[:])
}
