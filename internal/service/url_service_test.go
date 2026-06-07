package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/topinambur02/url-shortener/internal/dto"
	"github.com/topinambur02/url-shortener/internal/model"
	"github.com/topinambur02/url-shortener/internal/repository/mocks"
)

func TestURLService_GetByShortURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		shortUrl      string
		mockSetup     func(mockRepo *mocks.StoreRepository)
		expectedError error
		expectedRes   *dto.OriginalURLDto
	}{
		{
			name:     "Success",
			shortUrl: "abcd123",
			mockSetup: func(mockRepo *mocks.StoreRepository) {
				mockRepo.On("GetByShortURL", mock.Anything, "abcd123").
					Return(&model.URL{OriginalURL: "https://example.com", ShortURL: "abcd123"}, nil).
					Once()
			},
			expectedError: nil,
			expectedRes:   &dto.OriginalURLDto{OriginalURL: "https://example.com"},
		},
		{
			name:     "Error_NotFound",
			shortUrl: "notfound",
			mockSetup: func(mockRepo *mocks.StoreRepository) {
				mockRepo.On("GetByShortURL", mock.Anything, "notfound").
					Return((*model.URL)(nil), errors.New("url not found")).
					Once()
			},
			expectedError: errors.New("url not found"),
			expectedRes:   nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := mocks.NewStoreRepository(t)
			tt.mockSetup(mockRepo)
			service := NewURLService(mockRepo)

			res, err := service.GetByShortURL(context.Background(), tt.shortUrl)

			if tt.expectedError != nil {
				require.Error(t, err)
				require.EqualError(t, err, tt.expectedError.Error())
				require.Nil(t, res)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedRes, res)
			}
		})
	}
}

func TestURLService_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		createDto     *dto.CreateURLDto
		mockSetup     func(mockRepo *mocks.StoreRepository)
		expectedError error
		expectedRes   *dto.ShortURLDto
	}{
		{
			name:      "Success",
			createDto: &dto.CreateURLDto{OriginalURL: "https://example.com"},
			mockSetup: func(mockRepo *mocks.StoreRepository) {
				matcher := mock.MatchedBy(func(u model.URL) bool {
					return u.OriginalURL == "https://example.com"
				})

				mockRepo.On("Create", mock.Anything, matcher).
					Return(&model.URL{OriginalURL: "https://example.com", ShortURL: "genShort123"}, nil).
					Once()
			},
			expectedError: nil,
			expectedRes:   &dto.ShortURLDto{ShortURL: "http://localhost:8080/genShort123"},
		},
		{
			name:      "Error_RepositoryFails",
			createDto: &dto.CreateURLDto{OriginalURL: "https://broken.com"},
			mockSetup: func(mockRepo *mocks.StoreRepository) {
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("model.URL")).
					Return((*model.URL)(nil), errors.New("db error")).
					Once()
			},
			expectedError: errors.New("db error"),
			expectedRes:   nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := mocks.NewStoreRepository(t)
			tt.mockSetup(mockRepo)
			service := NewURLService(mockRepo)

			res, err := service.Create(context.Background(), tt.createDto, "localhost:8080")

			if tt.expectedError != nil {
				require.Error(t, err)
				require.EqualError(t, err, tt.expectedError.Error())
				require.Nil(t, res)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedRes, res)
			}
		})
	}
}
