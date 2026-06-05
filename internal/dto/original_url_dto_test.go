package dto

import (
	"encoding/json"
	"testing"
)

func TestOriginalUrlDto(t *testing.T) {
	t.Run("TestOriginalUrlDto_Unmarshal", func(t *testing.T) {
		tests := []struct {
			name      string
			inputJSON string
			expected  string
		}{
			{
				name:      "Valid JSON",
				inputJSON: `{"original_url":"https://example.com"}`,
				expected:  "https://example.com",
			},
			{
				name:      "Empty string value",
				inputJSON: `{"original_url":""}`,
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
				var dto OriginalUrlDto
				err := json.Unmarshal([]byte(tt.inputJSON), &dto)
				if err != nil {
					t.Fatalf("failed to unmarshal JSON: %v", err)
				}

				if dto.OriginalUrl != tt.expected {
					t.Errorf("got %q, want %q", dto.OriginalUrl, tt.expected)
				}
			})
		}
	})
	t.Run("TestOriginalUrlDto_Marshal", func(t *testing.T) {
		dto := OriginalUrlDto{
			OriginalUrl: "https://example.com",
		}
		expectedJSON := `{"original_url":"https://example.com"}`

		bytes, err := json.Marshal(dto)
		if err != nil {
			t.Fatalf("failed to marshal struct: %v", err)
		}

		if string(bytes) != expectedJSON {
			t.Errorf("got %s, want %s", string(bytes), expectedJSON)
		}
	})
}
