package dto

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorDto_Marshal(t *testing.T) {
	tests := []struct {
		name     string
		input    ErrorDto
		expected string
	}{
		{
			name:     "Successful serialization (404 Not Found)",
			input:    ErrorDto{StatusCode: 404, Message: "Not Found"},
			expected: `{"status_code":404,"message":"Not Found"}`,
		},
		{
			name:     "Serialization with empty values",
			input:    ErrorDto{StatusCode: 0, Message: ""},
			expected: `{"status_code":0,"message":""}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := json.Marshal(tt.input)
			require.NoError(t, err, "Failed to serialize structure")
			require.JSONEq(t, tt.expected, string(result), "The JSON should match what is expected.")
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
			name:     "Successful deserialization (500 Internal Server Error)",
			input:    `{"status_code":500,"message":"Internal Server Error"}`,
			expected: ErrorDto{StatusCode: 500, Message: "Internal Server Error"},
			wantErr:  false,
		},
		{
			name:    "Invalid JSON",
			input:   `{"status_code": 200, "message": `,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result ErrorDto
			err := json.Unmarshal([]byte(tt.input), &result)

			if tt.wantErr {
				require.Error(t, err, "deserialization error expected")
			} else {
				require.NoError(t, err, "deserialization should not cause an error")
				require.Equal(t, tt.expected, result, "the deserialized structure must match")
			}
		})
	}
}
