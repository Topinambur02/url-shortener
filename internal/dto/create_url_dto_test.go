package dto

import (
	"encoding/json"
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestCreateUrlDto(t *testing.T) {
	t.Run("TestCreateUrlDto_Validation", func(t *testing.T) {
		v := validator.New()

		tests := []struct {
			name    string
			dto     CreateUrlDto
			wantErr bool
		}{
			{
				name: "Valid URL",
				dto: CreateUrlDto{
					OriginalUrl: "https://github.com",
				},
				wantErr: false,
			},
			{
				name: "Empty URL (Required failure)",
				dto: CreateUrlDto{
					OriginalUrl: "",
				},
				wantErr: true,
			},
			{
				name: "Invalid URL format",
				dto: CreateUrlDto{
					OriginalUrl: "not-a-url",
				},
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := v.Struct(tt.dto)
				if (err != nil) != tt.wantErr {
					t.Errorf("v.Struct() error = %v, wantErr %v", err, tt.wantErr)
				}
			})
		}
	})
	t.Run("TestCreateUrlDto_Marshal", func(t *testing.T) {
		tests := []struct {
			name    string
			dto     CreateUrlDto
			want    string
			wantErr bool
		}{
			{
				name: "Valid struct to JSON",
				dto: CreateUrlDto{
					OriginalUrl: "https://example.com",
				},
				want:    `{"original_url":"https://example.com"}`,
				wantErr: false,
			},
			{
				name:    "Empty field struct to JSON",
				dto:     CreateUrlDto{OriginalUrl: ""},
				want:    `{"original_url":""}`,
				wantErr: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := json.Marshal(tt.dto)
				if (err != nil) != tt.wantErr {
					t.Errorf("json.Marshal() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !tt.wantErr && string(got) != tt.want {
					t.Errorf("json.Marshal() got = %s, want %s", string(got), tt.want)
				}
			})
		}
	})
	t.Run("TestCreateUrlDto_Unmarshal", func(t *testing.T) {
		tests := []struct {
			name    string
			jsonStr string
			want    CreateUrlDto
			wantErr bool
		}{
			{
				name:    "Valid JSON to struct",
				jsonStr: `{"original_url":"https://example.com"}`,
				want: CreateUrlDto{
					OriginalUrl: "https://example.com",
				},
				wantErr: false,
			},
			{
				name:    "JSON with extra fields (should ignore by default)",
				jsonStr: `{"original_url":"https://example.com","extra_field":"ignored"}`,
				want: CreateUrlDto{
					OriginalUrl: "https://example.com",
				},
				wantErr: false,
			},
			{
				name:    "Invalid JSON syntax",
				jsonStr: `{"original_url": "missing_quotes_and_brace`,
				want:    CreateUrlDto{},
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var got CreateUrlDto
				err := json.Unmarshal([]byte(tt.jsonStr), &got)
				if (err != nil) != tt.wantErr {
					t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !tt.wantErr && got != tt.want {
					t.Errorf("json.Unmarshal() got = %+v, want %+v", got, tt.want)
				}
			})
		}
	})
}
