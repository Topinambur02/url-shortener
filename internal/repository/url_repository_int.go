package repository

import (
	"context"

	"github.com/topinambur02/url-shortener/internal/model"
)

//go:generate mockery --name UrlRepository --output ./mocks --case underscore
type UrlRepository interface {
	GetByShortUrl(ctx context.Context, short_url string) (*model.Url, error)
	Create(ctx context.Context, url model.Url) (*model.Url, error)
}
