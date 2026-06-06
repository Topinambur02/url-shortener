package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/topinambur02/url-shortener/pkg/constants"
)

func TestGenerateShortUrl(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"Empty string", ""},
		{"Standard HTTP", "http://example.com"},
		{"Standard HTTPS", "https://google.com"},
		{"Long URL with params", "https://example.com/path/to/resource?id=12345&token=abcde"},
		{"URL with special characters", "https://example.com/path/!@#$%^&*()"},
		{"Non-ASCII characters", "https://пример.рф/тест"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateShortUrl(tt.input)

			require.Len(t, got, constants.Length, "expected length to be %d, got %q", constants.Length, got)

			for _, char := range got {
				require.Truef(t, strings.ContainsRune(constants.Alphabet, char), "GenerateShortUrl() contains invalid char: %c in string %q", char, got)
			}

			gotSecondTime := GenerateShortUrl(tt.input)
			require.Equal(t, got, gotSecondTime, "GenerateShortUrl() is not deterministic")
		})
	}

	t.Run("Different inputs give different outputs", func(t *testing.T) {
		url1 := "https://example.com/1"
		url2 := "https://example.com/2"

		hash1 := GenerateShortUrl(url1)
		hash2 := GenerateShortUrl(url2)

		require.NotEqual(t, hash1, hash2, "Collision detected for %q and %q", url1, url2)
	})
}