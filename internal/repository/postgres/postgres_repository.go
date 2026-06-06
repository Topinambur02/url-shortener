package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/topinambur02/url-shortener/internal/config"
	"github.com/topinambur02/url-shortener/internal/db"
	"github.com/topinambur02/url-shortener/internal/model"
	"github.com/topinambur02/url-shortener/internal/repository"
	"github.com/topinambur02/url-shortener/pkg/exceptions"
)

func init() {
	repository.Register("postgres", func(ctx context.Context, cfg *config.Config) (repository.StoreRepository, error) {
		database, err := db.InitDB(ctx, cfg)
		if err != nil {
			return nil, err
		}
		return NewPostgresRepository(database), nil
	})
}

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) repository.StoreRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) GetByShortURL(ctx context.Context, shortURL string) (*model.URL, error) {
	var url model.URL
	err := r.db.QueryRow(ctx,
		"SELECT id, short_url, original_url FROM urls WHERE short_url = $1",
		shortURL,
	).Scan(&url.ID, &url.ShortURL, &url.OriginalURL)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, exceptions.ErrNotFound
		}
		
		return nil, err
	}

	return &url, nil
}

func (r *PostgresRepository) Create(ctx context.Context, url model.URL) (*model.URL, error) {
	query := `INSERT INTO urls (original_url, short_url) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(ctx, query, url.OriginalURL, url.ShortURL)

	if err != nil {
		return nil, err
	}

	return &url, err
}
