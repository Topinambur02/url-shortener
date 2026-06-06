package dto

import (
	"encoding/json"
	"testing"
)

func TestErrorDto_Marshal(t *testing.T) {
	tests := []struct {
		name     string
		input    ErrorDto
		expected string
	}{
		{
			name: "Успешная сериализация (404 Not Found)",
			input: ErrorDto{
				StatusCode: 404,
				Message:    "Not Found",
			},
			expected: `{"status_code":404,"message":"Not Found"}`,
		},
		{
			name: "Сериализация с пустыми значениями",
			input: ErrorDto{
				StatusCode: 0,
				Message:    "",
			},
			expected: `{"status_code":0,"message":""}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("Не удалось сериализовать структуру: %v", err)
			}

			if string(result) != tt.expected {
				t.Errorf("Ожидалось:\n%s\nПолучено:\n%s", tt.expected, string(result))
			}
		})
	}
}

func TestErrorDto_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ErrorDto
		wantErr  bool
	}{
		{
			name:  "Успешная десериализация (500 Internal Server Error)",
			input: `{"status_code":500,"message":"Internal Server Error"}`,
			expected: ErrorDto{
				StatusCode: 500,
				Message:    "Internal Server Error",
			},
			wantErr: false,
		},
		{
			name:    "Невалидный JSON",
			input:   `{"status_code": 200, "message": `,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result ErrorDto
			err := json.Unmarshal([]byte(tt.input), &result)

			if (err != nil) != tt.wantErr {
				t.Fatalf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if result.StatusCode != tt.expected.StatusCode {
					t.Errorf("StatusCode: ожидалось %d, получено %d", tt.expected.StatusCode, result.StatusCode)
				}
				if result.Message != tt.expected.Message {
					t.Errorf("Message: ожидалось %q, получено %q", tt.expected.Message, result.Message)
				}
			}
		})
	}
}