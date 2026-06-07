package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/topinambur02/url-shortener/pkg/logging"
)

func setupLogging(t *testing.T) {
    logging.Init()

    t.Cleanup(func() {
        os.RemoveAll("logs")
    })
}

func TestLoggingResponseWriter(t *testing.T) {
	setupLogging(t)

	t.Run("capture status code and size", func(t *testing.T) {
		rec := httptest.NewRecorder()
		lrw := newLoggingResponseWriter(rec)

		require.Equal(t, http.StatusOK, lrw.statusCode)

		lrw.WriteHeader(http.StatusCreated)
		body := []byte("hello")
		n, err := lrw.Write(body)

		require.NoError(t, err)
		require.Equal(t, len(body), n)

		require.Equal(t, http.StatusCreated, lrw.statusCode)
		require.Equal(t, int64(len(body)), lrw.size)

		require.Equal(t, http.StatusCreated, rec.Code)
		require.Equal(t, "hello", rec.Body.String())
	})
}

func TestHTTPLoggerMiddleware(t *testing.T) {
	t.Parallel()
	setupLogging(t)

	tests := []struct {
		name           string
		requestID      string
		handlerStatus  int
		handlerPayload string
	}{
		{
			name:           "Success response (200 Info)",
			requestID:      "req-123",
			handlerStatus:  http.StatusOK,
			handlerPayload: "ok",
		},
		{
			name:           "Client error (400 Warn)",
			requestID:      "",
			handlerStatus:  http.StatusBadRequest,
			handlerPayload: "bad request",
		},
		{
			name:           "Server error (500 Error)",
			requestID:      "req-500",
			handlerStatus:  http.StatusInternalServerError,
			handlerPayload: "internal error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.handlerStatus)
				_, _ = w.Write([]byte(tt.handlerPayload))
			})

			middlewareToTest := HTTPLoggerMiddleware(nextHandler)

			req, err := http.NewRequest(http.MethodPost, "/api/v1/shorten", nil)
			require.NoError(t, err)

			if tt.requestID != "" {
				req.Header.Set("X-Request-ID", tt.requestID)
			}
			req.Header.Set("User-Agent", "Go-Test-Client")

			rec := httptest.NewRecorder()

			middlewareToTest.ServeHTTP(rec, req)

			require.Equal(t, tt.handlerStatus, rec.Code)
			require.Equal(t, tt.handlerPayload, rec.Body.String())
		})
	}
}
