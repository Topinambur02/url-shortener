package utils

import (
	"hash/fnv"

	"github.com/topinambur02/url-shortener/pkg/constants"
)

func GenerateShortUrl(originalUrl string) string {
	h := fnv.New64a()
	h.Write([]byte(originalUrl))
	hashValue := h.Sum64()
	var result [constants.Length]byte
	base := uint64(len(constants.Alphabet))

	for i := range constants.Length {
		result[i] = constants.Alphabet[hashValue % base]
		hashValue /= base
	}

	return string(result[:])
}
