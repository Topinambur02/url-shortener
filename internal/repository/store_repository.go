package repository

import (
	"context"

	"github.com/topinambur02/url-shortener/internal/model"
)

//go:generate mockery --name StoreRepository --output ./mocks --case underscore
type StoreRepository interface {
	GetByShortUrl(ctx context.Context, short_url string) (*model.URL, error)
	Create(ctx context.Context, url model.URL) (*model.URL, error)
}
