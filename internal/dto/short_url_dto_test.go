package dto

import (
	"encoding/json"
	"testing"
)

func TestShortUrlDto(t *testing.T) {
	t.Run("TestShortUrlDto_Unmarshal", func(t *testing.T) {
		tests := []struct {
			name      string
			inputJSON string
			expected  string
		}{
			{
				name:      "Valid JSON",
				inputJSON: `{"short_url":"https://sh.rt/abcde"}`,
				expected:  "https://sh.rt/abcde",
			},
			{
				name:      "Empty string value",
				inputJSON: `{"short_url":""}`,
				expected:  "",
			},
			{
				name:      "Missing field in JSON",
				inputJSON: `{}`,
				expected:  "",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var dto ShortUrlDto
				err := json.Unmarshal([]byte(tt.inputJSON), &dto)
				if err != nil {
					t.Fatalf("failed to unmarshal JSON: %v", err)
				}

				if dto.ShortUrl != tt.expected {
					t.Errorf("got %q, want %q", dto.ShortUrl, tt.expected)
				}
			})
		}
	})

	t.Run("TestShortUrlDto_Marshal", func(t *testing.T) {
		dto := ShortUrlDto{
			ShortUrl: "https://sh.rt/abcde",
		}
		expectedJSON := `{"short_url":"https://sh.rt/abcde"}`

		bytes, err := json.Marshal(dto)
		if err != nil {
			t.Fatalf("failed to marshal struct: %v", err)
		}

		if string(bytes) != expectedJSON {
			t.Errorf("got %s, want %s", string(bytes), expectedJSON)
		}
	})
}