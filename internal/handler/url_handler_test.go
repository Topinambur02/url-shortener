package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/topinambur02/url-shortener/internal/dto"
	"github.com/topinambur02/url-shortener/internal/service/mocks"
	"github.com/topinambur02/url-shortener/pkg/exceptions"
	"github.com/topinambur02/url-shortener/pkg/logging"
)

func setupLogging(t *testing.T) {
	logging.Init()

	t.Cleanup(func() {
		os.RemoveAll("logs")
	})
}

func TestURLHandler(t *testing.T) {
	setupLogging(t)
	t.Run("TestURLHandler_GetByShortUrl", func(t *testing.T) {
		tests := []struct {
			name           string
			shortURL       string
			mockBehavior   func(m *mocks.URLService)
			expectedStatus int
			expectedBody   string
		}{
			{
				name:           "Incorrect shortURL length (too short)",
				shortURL:       "abc",
				mockBehavior:   func(m *mocks.URLService) {},
				expectedStatus: http.StatusBadRequest,
				expectedBody:   "{\"status_code\":400,\"message\":\"invalid short url length\"}\n",
			},
			{
				name:     "Link not found (error 404)",
				shortURL: "1234567890",
				mockBehavior: func(m *mocks.URLService) {
					m.On("GetByShortURL", mock.Anything, "1234567890").
						Return((*dto.OriginalURLDto)(nil), exceptions.ErrNotFound)
				},
				expectedStatus: http.StatusNotFound,
				expectedBody:   "{\"status_code\":404,\"message\":\"not found\"}\n",
			},
			{
				name:     "Internal service error (error 500)",
				shortURL: "1234567890",
				mockBehavior: func(m *mocks.URLService) {
					m.On("GetByShortURL", mock.Anything, "1234567890").
						Return((*dto.OriginalURLDto)(nil), errors.New("db connection failure"))
				},
				expectedStatus: http.StatusInternalServerError,
				expectedBody:   "{\"status_code\":500,\"message\":\"internal error\"}\n",
			},
			{
				name:     "Successfully retrieved the original URL",
				shortURL: "1234567890",
				mockBehavior: func(m *mocks.URLService) {
					m.On("GetByShortURL", mock.Anything, "1234567890").
						Return(&dto.OriginalURLDto{OriginalURL: "https://example.com"}, nil)
				},
				expectedStatus: http.StatusOK,
				expectedBody:   `{"original_url":"https://example.com"}`,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockSvc := mocks.NewURLService(t)
				tt.mockBehavior(mockSvc)
				h := NewURLHandler(mockSvc, "http://localhost:8080")

				req := httptest.NewRequest(http.MethodGet, "/"+tt.shortURL, nil)
				req.SetPathValue("short", tt.shortURL)

				rr := httptest.NewRecorder()
				h.GetByShortURL(rr, req)

				require.Equal(t, tt.expectedStatus, rr.Code)

				if tt.expectedStatus == http.StatusOK {
					require.Contains(t, rr.Body.String(), "https://example.com")
				} else {
					require.Equal(t, tt.expectedBody, rr.Body.String())
				}
			})
		}
	})

	t.Run("TestURLHandler_Create", func(t *testing.T) {
		tests := []struct {
			name           string
			reqBody        string
			mockBehavior   func(m *mocks.URLService)
			expectedStatus int
			checkBody      func(t *testing.T, body string)
		}{
			{
				name:           "Invalid JSON in request",
				reqBody:        `{invalid json`,
				mockBehavior:   func(m *mocks.URLService) {},
				expectedStatus: http.StatusBadRequest,
				checkBody: func(t *testing.T, body string) {
					require.Contains(t, body, "invalid request")
				},
			},
			{
				name:           "DTO validation error (empty string passed)",
				reqBody:        `{"original_url": ""}`,
				mockBehavior:   func(m *mocks.URLService) {},
				expectedStatus: http.StatusBadRequest,
				checkBody: func(t *testing.T, body string) {
					require.Contains(t, body, "Validation failed")
				},
			},
			{
				name:    "Internal error creating link",
				reqBody: `{"original_url": "https://example.com"}`,
				mockBehavior: func(m *mocks.URLService) {
					m.On("Create", mock.Anything, &dto.CreateURLDto{OriginalURL: "https://example.com"}, "http://localhost:8080").
						Return((*dto.ShortURLDto)(nil), errors.New("failed to insert"))
				},
				expectedStatus: http.StatusInternalServerError,
				checkBody: func(t *testing.T, body string) {
					require.Contains(t, body, "internal error")
				},
			},
			{
				name:    "Successful creation of a short link",
				reqBody: `{"original_url": "https://example.com"}`,
				mockBehavior: func(m *mocks.URLService) {
					m.On("Create", mock.Anything, &dto.CreateURLDto{OriginalURL: "https://example.com"}, "http://localhost:8080").
						Return(&dto.ShortURLDto{ShortURL: "1234567890"}, nil)
				},
				expectedStatus: http.StatusCreated,
				checkBody: func(t *testing.T, body string) {
					require.Contains(t, body, "1234567890")
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockSvc := mocks.NewURLService(t)
				tt.mockBehavior(mockSvc)
				h := NewURLHandler(mockSvc, "http://localhost:8080")

				req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tt.reqBody))
				req.Header.Set("Content-Type", "application/json")

				rr := httptest.NewRecorder()
				h.Create(rr, req)

				require.Equal(t, tt.expectedStatus, rr.Code)
				tt.checkBody(t, rr.Body.String())
			})
		}
	})
}
