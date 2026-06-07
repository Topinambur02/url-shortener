package dto

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShortURLDto(t *testing.T) {
	t.Run("Unmarshal", func(t *testing.T) {
		tests := []struct {
			name      string
			inputJSON string
			expected  string
		}{
			{"Valid JSON", `{"short_url":"https://sh.rt/abcde"}`, "https://sh.rt/abcde"},
			{"Empty string value", `{"short_url":""}`, ""},
			{"Missing field in JSON", `{}`, ""},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var dto ShortURLDto
				err := json.Unmarshal([]byte(tt.inputJSON), &dto)
				require.NoError(t, err, "Unmarshal should not throw an error")
				require.Equal(t, tt.expected, dto.ShortURL, "The ShortUrl field must match")
			})
		}
	})

	t.Run("Marshal", func(t *testing.T) {
		dto := ShortURLDto{ShortURL: "https://sh.rt/abcde"}
		expectedJSON := `{"short_url":"https://sh.rt/abcde"}`

		bytes, err := json.Marshal(dto)
		require.NoError(t, err, "Marshal should not throw an error")
		require.JSONEq(t, expectedJSON, string(bytes), "JSON must match")
	})
}
