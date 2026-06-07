package dto

import (
	"encoding/json"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestCreateURLDto(t *testing.T) {
	t.Run("Validation", func(t *testing.T) {
		v := validator.New()
		tests := []struct {
			name    string
			dto     CreateURLDto
			wantErr bool
		}{
			{"Valid URL", CreateURLDto{OriginalURL: "https://github.com"}, false},
			{"Empty URL (Required failure)", CreateURLDto{OriginalURL: ""}, true},
			{"Invalid URL format", CreateURLDto{OriginalURL: "not-a-url"}, true},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := v.Struct(tt.dto)
				if tt.wantErr {
					require.Error(t, err, "a validation error was expected")
				} else {
					require.NoError(t, err, "no validation error expected")
				}
			})
		}
	})

	t.Run("Marshal", func(t *testing.T) {
		tests := []struct {
			name string
			dto  CreateURLDto
			want string
		}{
			{
				name: "Valid struct to JSON",
				dto:  CreateURLDto{OriginalURL: "https://example.com"},
				want: `{"original_url":"https://example.com"}`,
			},
			{
				name: "Empty field struct to JSON",
				dto:  CreateURLDto{OriginalURL: ""},
				want: `{"original_url":""}`,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := json.Marshal(tt.dto)
				require.NoError(t, err, "json.Marshal should not cause an error")
				require.JSONEq(t, tt.want, string(got), "The Marshal result must match the expected JSON")
			})
		}
	})

	t.Run("Unmarshal", func(t *testing.T) {
		tests := []struct {
			name    string
			jsonStr string
			want    CreateURLDto
			wantErr bool
		}{
			{
				name:    "Valid JSON to struct",
				jsonStr: `{"original_url":"https://example.com"}`,
				want:    CreateURLDto{OriginalURL: "https://example.com"},
				wantErr: false,
			},
			{
				name:    "JSON with extra fields (should ignore by default)",
				jsonStr: `{"original_url":"https://example.com","extra_field":"ignored"}`,
				want:    CreateURLDto{OriginalURL: "https://example.com"},
				wantErr: false,
			},
			{
				name:    "Invalid JSON syntax",
				jsonStr: `{"original_url": "missing_quotes_and_brace`,
				want:    CreateURLDto{},
				wantErr: true,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var got CreateURLDto
				err := json.Unmarshal([]byte(tt.jsonStr), &got)
				if tt.wantErr {
					require.Error(t, err, "An error was expected while Unmarshal")
				} else {
					require.NoError(t, err, "Unmarshal should not throw an error")
					require.Equal(t, tt.want, got, "the result of Unmarshal must match the expected structure")
				}
			})
		}
	})
}
