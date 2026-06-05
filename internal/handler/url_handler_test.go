package handler

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/topinambur02/url-shortener/internal/dto"
	"github.com/topinambur02/url-shortener/internal/service/mocks"
	"github.com/topinambur02/url-shortener/pkg/exceptions"
	"github.com/topinambur02/url-shortener/pkg/logging"
)

func TestMain(m *testing.M) {
	logging.Init()
	exitCode := m.Run()

	err := os.RemoveAll("logs")
	if err != nil {
		fmt.Printf("Log error when deleting folders: %v\n", err)
	}

	os.Exit(exitCode)
}

func TestUrlHandler(t *testing.T) {
	t.Run("TestURLHandler_GetByShortUrl", func(t *testing.T) {
		tests := []struct {
			name           string
			shortURL       string
			mockBehavior   func(m *mocks.UrlService)
			expectedStatus int
			expectedBody   string
		}{
			{
				name:           "Incorrect shortURL length (too short)",
				shortURL:       "abc",
				mockBehavior:   func(m *mocks.UrlService) {},
				expectedStatus: http.StatusBadRequest,
				expectedBody:   "invalid short url length\n",
			},
			{
				name:     "Link not found (error 404)",
				shortURL: "1234567890",
				mockBehavior: func(m *mocks.UrlService) {
					m.On("GetByShortUrl", mock.Anything, "1234567890").
						Return((*dto.OriginalUrlDto)(nil), exceptions.ErrNotFound)
				},
				expectedStatus: http.StatusNotFound,
				expectedBody:   "not found\n",
			},
			{
				name:     "Internal service error (error 500)",
				shortURL: "1234567890",
				mockBehavior: func(m *mocks.UrlService) {
					m.On("GetByShortUrl", mock.Anything, "1234567890").
						Return((*dto.OriginalUrlDto)(nil), errors.New("db connection failure"))
				},
				expectedStatus: http.StatusInternalServerError,
				expectedBody:   "internal error\n",
			},
			{
				name:     "Successfully retrieved the original URL",
				shortURL: "1234567890",
				mockBehavior: func(m *mocks.UrlService) {
					m.On("GetByShortUrl", mock.Anything, "1234567890").
						Return(&dto.OriginalUrlDto{OriginalUrl: "https://example.com"}, nil)
				},
				expectedStatus: http.StatusOK,
				expectedBody:   `{"original_url":"https://example.com"}`,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockSvc := mocks.NewUrlService(t)
				tt.mockBehavior(mockSvc)
				h := NewURLHandler(mockSvc)

				req := httptest.NewRequest(http.MethodGet, "/"+tt.shortURL, nil)
				req.SetPathValue("short", tt.shortURL)

				rr := httptest.NewRecorder()
				h.GetByShortUrl(rr, req)

				if rr.Code != tt.expectedStatus {
					t.Errorf("expected status %d, received %d", tt.expectedStatus, rr.Code)
				}

				if tt.expectedStatus == http.StatusOK {
					if !strings.Contains(rr.Body.String(), "https://example.com") {
						t.Errorf("A valid JSON response was expected, received: %q", rr.Body.String())
					}
				} else {
					if rr.Body.String() != tt.expectedBody {
						t.Errorf("expected body %q, received %q", tt.expectedBody, rr.Body.String())
					}
				}
			})
		}
	})

	t.Run("TestURLHandler_Create", func(t *testing.T) {
		tests := []struct {
			name           string
			reqBody        string
			mockBehavior   func(m *mocks.UrlService)
			expectedStatus int
			checkBody      func(t *testing.T, body string)
		}{
			{
				name:           "Invalid JSON in request",
				reqBody:        `{invalid json`,
				mockBehavior:   func(m *mocks.UrlService) {},
				expectedStatus: http.StatusBadRequest,
				checkBody: func(t *testing.T, body string) {
					if !strings.Contains(body, "invalid request") {
						t.Errorf("expected text 'invalid request', received %q", body)
					}
				},
			},
			{
				name:           "DTO validation error (empty string passed)",
				reqBody:        `{"original_url": ""}`,
				mockBehavior:   func(m *mocks.UrlService) {},
				expectedStatus: http.StatusBadRequest,
				checkBody: func(t *testing.T, body string) {
					if !strings.Contains(body, "Validation failed") {
						t.Errorf("expected text 'Validation failed', received %q", body)
					}
				},
			},
			{
				name:    "Internal error creating link",
				reqBody: `{"original_url": "https://example.com"}`,
				mockBehavior: func(m *mocks.UrlService) {
					m.On("Create", mock.Anything, &dto.CreateUrlDto{OriginalUrl: "https://example.com"}).
						Return((*dto.ShortUrlDto)(nil), errors.New("failed to insert"))
				},
				expectedStatus: http.StatusInternalServerError,
				checkBody: func(t *testing.T, body string) {
					if !strings.Contains(body, "internal error") {
						t.Errorf("expected text 'internal error', received %q", body)
					}
				},
			},
			{
				name:    "Successful creation of a short link",
				reqBody: `{"original_url": "https://example.com"}`,
				mockBehavior: func(m *mocks.UrlService) {
					m.On("Create", mock.Anything, &dto.CreateUrlDto{OriginalUrl: "https://example.com"}).
						Return(&dto.ShortUrlDto{ShortUrl: "1234567890"}, nil)
				},
				expectedStatus: http.StatusOK,
				checkBody: func(t *testing.T, body string) {
					if !strings.Contains(body, "1234567890") {
						t.Errorf("expected short URL in response, received %q", body)
					}
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockSvc := mocks.NewUrlService(t)
				tt.mockBehavior(mockSvc)
				h := NewURLHandler(mockSvc)

				req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tt.reqBody))
				req.Header.Set("Content-Type", "application/json")

				rr := httptest.NewRecorder()
				h.Create(rr, req)

				if rr.Code != tt.expectedStatus {
					t.Errorf("expected status %d, received %d", tt.expectedStatus, rr.Code)
				}

				tt.checkBody(t, rr.Body.String())
			})
		}
	})
}
