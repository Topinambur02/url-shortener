package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/topinambur02/url-shortener/pkg/logging"
)

func TestMain(m *testing.M) {
	logging.Init()
	exitCode := m.Run()

	err := os.RemoveAll("logs")
	if err != nil {
		fmt.Printf("Error deleting logs folder: %v\n", err)
	}

	os.Exit(exitCode)
}

func TestLoggingResponseWriter(t *testing.T) {
	t.Run("capture status code and size", func(t *testing.T) {
		rec := httptest.NewRecorder()
		lrw := newLoggingResponseWriter(rec)

		if lrw.statusCode != http.StatusOK {
			t.Errorf("expected default status OK (200), got %d", lrw.statusCode)
		}

		lrw.WriteHeader(http.StatusCreated)
		body := []byte("hello")
		n, err := lrw.Write(body)

		if err != nil {
			t.Fatalf("failed to write body: %v", err)
		}
		if n != len(body) {
			t.Errorf("expected written bytes %d, got %d", len(body), n)
		}

		if lrw.statusCode != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, lrw.statusCode)
		}
		if lrw.size != int64(len(body)) {
			t.Errorf("expected size %d, got %d", len(body), lrw.size)
		}

		if rec.Code != http.StatusCreated {
			t.Errorf("expected recorded status %d, got %d", http.StatusCreated, rec.Code)
		}
		if rec.Body.String() != "hello" {
			t.Errorf("expected recorded body 'hello', got '%s'", rec.Body.String())
		}
	})
}

func TestHTTPLoggerMiddleware(t *testing.T) {
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
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			
			if tt.requestID != "" {
				req.Header.Set("X-Request-ID", tt.requestID)
			}
			req.Header.Set("User-Agent", "Go-Test-Client")

			rec := httptest.NewRecorder()

			middlewareToTest.ServeHTTP(rec, req)

			if rec.Code != tt.handlerStatus {
				t.Errorf("expected status %d, got %d", tt.handlerStatus, rec.Code)
			}
			if rec.Body.String() != tt.handlerPayload {
				t.Errorf("expected body %q, got %q", tt.handlerPayload, rec.Body.String())
			}
		})
	}
}