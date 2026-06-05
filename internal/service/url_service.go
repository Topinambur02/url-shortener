package service

import (
	"context"

	"github.com/topinambur02/url-shortener/internal/dto"
	"github.com/topinambur02/url-shortener/internal/model"
	"github.com/topinambur02/url-shortener/internal/repository"
	"github.com/topinambur02/url-shortener/pkg/utils"
)

type UrlService interface {
	GetByShortUrl(ctx context.Context, shortUrl string) (*dto.OriginalUrlDto, error)
	Create(ctx context.Context, createUrlDto *dto.CreateUrlDto) (*dto.ShortUrlDto, error)
}

type UrlServiceImpl struct {
	urlRepo repository.UrlRepository
}

func NewUrlService(urlRepo repository.UrlRepository) *UrlServiceImpl {
	return &UrlServiceImpl{
		urlRepo: urlRepo,
	}
}

func (s *UrlServiceImpl) GetByShortUrl(ctx context.Context, shortUrl string) (*dto.OriginalUrlDto, error) {
	url, err := s.urlRepo.GetByShortUrl(ctx, shortUrl)

	if err != nil {
		return nil, err
	}

	return &dto.OriginalUrlDto{OriginalUrl: url.OriginalUrl}, nil
}

func (s *UrlServiceImpl) Create(ctx context.Context, createUrlDto *dto.CreateUrlDto) (*dto.ShortUrlDto, error) {
	originalUrl := createUrlDto.OriginalUrl
	shortURL := utils.GenerateShortUrl(originalUrl)
	url_hash := utils.HashURL(shortURL)

	url := model.Url{
		OriginalUrl: originalUrl,
		ShortUrl:    shortURL,
		UrlHash:     url_hash,
	}

	createdUrl, err := s.urlRepo.Create(ctx, url)

	if err != nil {
		return nil, err
	}

	return &dto.ShortUrlDto{ShortUrl: createdUrl.ShortUrl}, nil
}
