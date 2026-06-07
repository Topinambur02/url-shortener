package exceptions

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
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

			require.Equal(t, tt.statusCode, rec.Code, "HTTP status code mismatch")
			require.Equal(t, "application/json", rec.Header().Get("Content-Type"), "Content-Type header mismatch")

			var actualDto dto.ErrorDto
			err := json.NewDecoder(rec.Body).Decode(&actualDto)
			require.NoError(t, err, "failed to decode response body")

			expectedDto := dto.ErrorDto{
				StatusCode: tt.statusCode,
				Message:    tt.message,
			}
			require.Equal(t, expectedDto, actualDto, "JSON response body mismatch")
		})
	}
}
