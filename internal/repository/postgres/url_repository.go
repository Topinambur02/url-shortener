package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/topinambur02/url-shortener/internal/model"
	"github.com/topinambur02/url-shortener/internal/repository"
	"github.com/topinambur02/url-shortener/pkg/exceptions"
)

type UrlRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewUrlRepository(db *pgxpool.Pool) repository.UrlRepository {
	return &UrlRepositoryImpl{
		db: db,
	}
}

func (r *UrlRepositoryImpl) GetByShortUrl(ctx context.Context, short_url string) (*model.Url, error) {
	var url model.Url
	err := r.db.QueryRow(ctx,
		"SELECT id, short_url, original_url, url_hash FROM urls WHERE short_url = $1",
		short_url,
	).Scan(&url.ID, &url.ShortUrl, &url.OriginalUrl, &url.UrlHash)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, exceptions.ErrNotFound
		}
		
		return nil, err
	}

	return &url, nil
}

func (r *UrlRepositoryImpl) Create(ctx context.Context, url model.Url) (*model.Url, error) {
	query := `INSERT INTO urls (original_url, short_url, url_hash) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(ctx, query, url.OriginalUrl, url.ShortUrl, url.UrlHash)

	if err != nil {
		return nil, err
	}

	return &url, err
}
