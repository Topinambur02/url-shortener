package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/topinambur02/url-shortener/internal/model"
	"github.com/topinambur02/url-shortener/internal/repository"
)

type UrlRepositoryImpl struct {
	db *sql.DB
}

func NewUrlRepository(db *sql.DB) repository.UrlRepository {
	return &UrlRepositoryImpl{
		db: db,
	}
}

func (r *UrlRepositoryImpl) GetByShortUrl(ctx context.Context, short_url string) (*model.Url, error) {
	var url model.Url
	err := r.db.QueryRowContext(ctx, 
        "SELECT id, short_url, original_url, url_hash FROM urls WHERE short_url = $1", 
        short_url,
    ).Scan(&url.ID, &url.ShortUrl, &url.OriginalUrl, &url.UrlHash)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return &url, err
}

func (r *UrlRepositoryImpl) Create(ctx context.Context, url model.Url) (*model.Url, error) {
	query := `INSERT INTO urls (original_url, short_url, url_hash) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, url.OriginalUrl, url.ShortUrl, url.UrlHash)

	if err != nil {
		return nil, err
	}

	return &url, err
}