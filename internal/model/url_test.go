package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUrl(t *testing.T) {
	t.Run("TestUrl_Marshal", func(t *testing.T) {
		tests := []struct {
			name     string
			input    URL
			expected string
		}{
			{
				name: "Successful serialization of a filled structure",
				input: URL{
					ID:          1,
					OriginalURL: "https://example.com/some/very/long/path",
					ShortURL:    "http://sh.rt/aBcd12",
				},
				expected: `{"id":1,"original_url":"https://example.com/some/very/long/path","short_url":"http://sh.rt/aBcd12"}`,
			},
			{
				name:     "Serializing an empty structure (default values)",
				input:    URL{},
				expected: `{"id":0,"original_url":"","short_url":""}`,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				bytes, err := json.Marshal(tt.input)
				require.NoError(t, err)
				require.Equal(t, tt.expected, string(bytes))
			})
		}
	})

	t.Run("TestUrl_Unmarshal", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected URL
			wantErr  bool
		}{
			{
				name:  "Successfully deserializing valid JSON",
				input: `{"id":42,"original_url":"https://google.com","url_hash":"xyz987","short_url":"http://sh.rt/xyz987"}`,
				expected: URL{
					ID:          42,
					OriginalURL: "https://google.com",
					ShortURL:    "http://sh.rt/xyz987",
				},
				wantErr: false,
			},
			{
				name:    "Error with invalid data type (id passed as string instead of uint)",
				input:   `{"id":"not_a_number","original_url":"https://google.com"}`,
				wantErr: true,
			},
			{
				name:    "Error due to broken JSON syntax",
				input:   `{"id": 1, "original_url": `,
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var result URL
				err := json.Unmarshal([]byte(tt.input), &result)

				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					require.Equal(t, tt.expected, result)
				}
			})
		}
	})
}
