package service

import (
	"context"

	"github.com/topinambur02/url-shortener/internal/dto"
	"github.com/topinambur02/url-shortener/internal/model"
	"github.com/topinambur02/url-shortener/internal/repository"
	"github.com/topinambur02/url-shortener/pkg/utils"
)

//go:generate mockery --name URLService --output ./mocks --case underscore
type URLService interface {
	GetByShortURL(ctx context.Context, shortURL string) (*dto.OriginalURLDto, error)
	Create(ctx context.Context, createURLDto *dto.CreateURLDto, address string) (*dto.ShortURLDto, error)
}

type URLServiceImpl struct {
	repo repository.StoreRepository
}

func NewURLService(repo repository.StoreRepository) *URLServiceImpl {
	return &URLServiceImpl{
		repo: repo,
	}
}

func (s *URLServiceImpl) GetByShortURL(ctx context.Context, shortURL string) (*dto.OriginalURLDto, error) {
	url, err := s.repo.GetByShortURL(ctx, shortURL)

	if err != nil {
		return nil, err
	}

	return &dto.OriginalURLDto{OriginalURL: url.OriginalURL}, nil
}

func (s *URLServiceImpl) Create(ctx context.Context, createURLDto *dto.CreateURLDto, address string) (*dto.ShortURLDto, error) {
	originalURL := createURLDto.OriginalURL
	shortURL := utils.GenerateShortUrl(originalURL)

	url := model.URL{
		OriginalURL: originalURL,
		ShortURL:    shortURL,
	}

	createdURL, err := s.repo.Create(ctx, url)

	if err != nil {
		return nil, err
	}

	link := "http://" + address + "/" + createdURL.ShortURL

	return &dto.ShortURLDto{ShortURL: link}, nil
}
