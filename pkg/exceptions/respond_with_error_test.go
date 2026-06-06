package exceptions

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/topinambur02/url-shortener/internal/dto"
)

func TestRespondWithError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		message    string
	}{
		{
			name:       "400 Bad Request",
			statusCode: http.StatusBadRequest,
			message:    "Invalid request payload",
		},
		{
			name:       "404 Not Found",
			statusCode: http.StatusNotFound,
			message:    "URL not found",
		},
		{
			name:       "500 Internal Server Error",
			statusCode: http.StatusInternalServerError,
			message:    "Something went wrong on our side",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			RespondWithError(rec, tt.statusCode, tt.message)

			if rec.Code != tt.statusCode {
				t.Errorf("expected status %d, got %d", tt.statusCode, rec.Code)
			}

			contentType := rec.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("expected Content-Type 'application/json', got '%s'", contentType)
			}

			var actualDto dto.ErrorDto
			err := json.NewDecoder(rec.Body).Decode(&actualDto)
			if err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}

			if actualDto.StatusCode != tt.statusCode {
				t.Errorf("expected JSON StatusCode %d, got %d", tt.statusCode, actualDto.StatusCode)
			}

			if actualDto.Message != tt.message {
				t.Errorf("expected JSON Message %q, got %q", tt.message, actualDto.Message)
			}
		})
	}
}