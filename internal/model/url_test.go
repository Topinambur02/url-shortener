package model

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestUrl(t *testing.T) {
	t.Run("TestUrl_Marshal", func(t *testing.T) {
		tests := []struct {
			name     string
			input    Url
			expected string
		}{
			{
				name: "Successful serialization of a filled structure",
				input: Url{
					ID:          1,
					OriginalUrl: "https://example.com/some/very/long/path",
					UrlHash:     "aBcd12",
					ShortUrl:    "http://sh.rt/aBcd12",
				},
				expected: `{"id":1,"original_url":"https://example.com/some/very/long/path","url_hash":"aBcd12","short_url":"http://sh.rt/aBcd12"}`,
			},
			{
				name:     "Serializing an empty structure (default values)",
				input:    Url{},
				expected: `{"id":0,"original_url":"","url_hash":"","short_url":""}`,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				bytes, err := json.Marshal(tt.input)
				if err != nil {
					t.Fatalf("Error while marshaling: %v", err)
				}

				if string(bytes) != tt.expected {
					t.Errorf("\nResult: %s\nExpected: %s", string(bytes), tt.expected)
				}
			})
		}
	})
	t.Run("TestUrl_Unmarshal", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected Url
			wantErr  bool
		}{
			{
				name:  "Successfully deserializing valid JSON",
				input: `{"id":42,"original_url":"https://google.com","url_hash":"xyz987","short_url":"http://sh.rt/xyz987"}`,
				expected: Url{
					ID:          42,
					OriginalUrl: "https://google.com",
					UrlHash:     "xyz987",
					ShortUrl:    "http://sh.rt/xyz987",
				},
				wantErr: false,
			},
			{
				name:     "Error with invalid data type (id passed as string instead of uint)",
				input:    `{"id":"not_a_number","original_url":"https://google.com"}`,
				expected: Url{},
				wantErr:  true,
			},
			{
				name:     "Error due to broken JSON syntax",
				input:    `{"id": 1, "original_url": `,
				expected: Url{},
				wantErr:  true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var result Url
				err := json.Unmarshal([]byte(tt.input), &result)

				if (err != nil) != tt.wantErr {
					t.Fatalf("Expected to have error = %v, but got error: %v", tt.wantErr, err)
				}

				if !tt.wantErr && !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("\nResult:  %+v\nExpected: %+v", result, tt.expected)
				}
			})
		}
	})
}
