package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/topinambur02/url-shortener/internal/dto"
	"github.com/topinambur02/url-shortener/internal/model"
	"github.com/topinambur02/url-shortener/internal/repository/mocks"
)

func TestUrlService(t *testing.T) {
	t.Run("TestUrlServiceImpl_GetByShortUrl", func(t *testing.T) {
		tests := []struct {
			name          string
			shortUrl      string
			mockSetup     func(mockRepo *mocks.UrlRepository)
			expectedError error
			expectedRes   *dto.OriginalUrlDto
		}{
			{
				name:     "Success",
				shortUrl: "abcd123",
				mockSetup: func(mockRepo *mocks.UrlRepository) {
					mockRepo.On("GetByShortUrl", mock.Anything, "abcd123").
						Return(&model.Url{OriginalUrl: "https://example.com", ShortUrl: "abcd123"}, nil).
						Once()
				},
				expectedError: nil,
				expectedRes:   &dto.OriginalUrlDto{OriginalUrl: "https://example.com"},
			},
			{
				name:     "Error_NotFound",
				shortUrl: "notfound",
				mockSetup: func(mockRepo *mocks.UrlRepository) {
					mockRepo.On("GetByShortUrl", mock.Anything, "notfound").
						Return((*model.Url)(nil), errors.New("url not found")).
						Once()
				},
				expectedError: errors.New("url not found"),
				expectedRes:   nil,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockRepo := mocks.NewUrlRepository(t)
				tt.mockSetup(mockRepo)
				service := NewUrlService(mockRepo)

				res, err := service.GetByShortUrl(context.Background(), tt.shortUrl)

				if tt.expectedError != nil {
					assert.Error(t, err)
					assert.Equal(t, tt.expectedError.Error(), err.Error())
					assert.Nil(t, res)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tt.expectedRes, res)
				}
			})
		}
	})
	t.Run("TestUrlServiceImpl_Create", func(t *testing.T) {
		tests := []struct {
			name          string
			createDto     *dto.CreateUrlDto
			mockSetup     func(mockRepo *mocks.UrlRepository)
			expectedError error
			expectedRes   *dto.ShortUrlDto
		}{
			{
				name:      "Success",
				createDto: &dto.CreateUrlDto{OriginalUrl: "https://example.com"},
				mockSetup: func(mockRepo *mocks.UrlRepository) {
					matcher := mock.MatchedBy(func(u model.Url) bool {
						return u.OriginalUrl == "https://example.com" && u.ShortUrl != ""
					})

					mockRepo.On("Create", mock.Anything, matcher).
						Return(&model.Url{OriginalUrl: "https://example.com", ShortUrl: "genShort123"}, nil).
						Once()
				},
				expectedError: nil,
				expectedRes:   &dto.ShortUrlDto{ShortUrl: "genShort123"},
			},
			{
				name:      "Error_RepositoryFails",
				createDto: &dto.CreateUrlDto{OriginalUrl: "https://broken.com"},
				mockSetup: func(mockRepo *mocks.UrlRepository) {
					matcher := mock.MatchedBy(func(u model.Url) bool {
						return u.OriginalUrl == "https://broken.com"
					})

					mockRepo.On("Create", mock.Anything, matcher).
						Return((*model.Url)(nil), errors.New("db error")).
						Once()
				},
				expectedError: errors.New("db error"),
				expectedRes:   nil,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockRepo := mocks.NewUrlRepository(t)
				tt.mockSetup(mockRepo)
				service := NewUrlService(mockRepo)

				res, err := service.Create(context.Background(), tt.createDto)

				if tt.expectedError != nil {
					assert.Error(t, err)
					assert.Equal(t, tt.expectedError.Error(), err.Error())
					assert.Nil(t, res)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tt.expectedRes, res)
				}
			})
		}
	})
}
