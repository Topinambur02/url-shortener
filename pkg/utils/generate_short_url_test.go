package utils

import (
	"strings"
	"testing"
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

			if len(got) != length {
				t.Errorf("GenerateShortUrl() length = %d, want %d. Got string: %q", len(got), length, got)
			}

			for _, char := range got {
				if !strings.ContainsRune(alphabet, char) {
					t.Errorf("GenerateShortUrl() contains invalid char: %c in string %q", char, got)
				}
			}

			gotSecondTime := GenerateShortUrl(tt.input)
			if got != gotSecondTime {
				t.Errorf("GenerateShortUrl() is not deterministic. First: %q, Second: %q", got, gotSecondTime)
			}
		})
	}

	t.Run("Different inputs give different outputs", func(t *testing.T) {
		url1 := "https://example.com/1"
		url2 := "https://example.com/2"

		hash1 := GenerateShortUrl(url1)
		hash2 := GenerateShortUrl(url2)

		if hash1 == hash2 {
			t.Errorf("Collision detected for %q and %q. Both gave %q", url1, url2, hash1)
		}
	})
}