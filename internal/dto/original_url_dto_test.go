package dto

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOriginalURLDto(t *testing.T) {
	t.Run("Unmarshal", func(t *testing.T) {
		tests := []struct {
			name      string
			inputJSON string
			expected  string
		}{
			{"Valid JSON", `{"original_url":"https://example.com"}`, "https://example.com"},
			{"Empty string value", `{"original_url":""}`, ""},
			{"Missing field in JSON", `{}`, ""},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var dto OriginalURLDto
				err := json.Unmarshal([]byte(tt.inputJSON), &dto)
				require.NoError(t, err, "Unmarshal should not throw an error")
				require.Equal(t, tt.expected, dto.OriginalURL, "The OriginalUrl field must match")
			})
		}
	})

	t.Run("Marshal", func(t *testing.T) {
		dto := OriginalURLDto{OriginalURL: "https://example.com"}
		expectedJSON := `{"original_url":"https://example.com"}`

		bytes, err := json.Marshal(dto)
		require.NoError(t, err, "Marshal should not throw an error")
		require.JSONEq(t, expectedJSON, string(bytes), "JSON must match")
	})
}
