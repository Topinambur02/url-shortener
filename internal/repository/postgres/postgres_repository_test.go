package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/topinambur02/url-shortener/internal/model"
	"github.com/topinambur02/url-shortener/pkg/exceptions"
)

func setupTestDB(t *testing.T) (*pgxpool.Pool, func()) {
	ctx := context.Background()

	postgresContainer, err := tcpostgres.Run(ctx,
		"postgres:15-alpine",
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("testuser"),
		tcpostgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	require.NoError(t, err)

	connString, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, connString)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		CREATE TABLE urls (
			id SERIAL PRIMARY KEY,
			short_url VARCHAR(255) UNIQUE NOT NULL,
			original_url TEXT NOT NULL
		);
	`)
	require.NoError(t, err)

	cleanup := func() {
		pool.Close()
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}

	return pool, cleanup
}

func TestUrlRepository_CreateAndGet(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewPostgresRepository(pool)
	ctx := context.Background()

	inputURL := model.URL{
		OriginalURL: "https://example.com",
		ShortURL:    "exmpl123",
	}
	t.Run("Create Success", func(t *testing.T) {
		created, err := repo.Create(ctx, inputURL)
		require.NoError(t, err)
		require.NotNil(t, created)
		require.Equal(t, inputURL.OriginalURL, created.OriginalURL)
	})
	t.Run("Get Success", func(t *testing.T) {
		found, err := repo.GetByShortURL(ctx, inputURL.ShortURL)
		require.NoError(t, err)
		require.NotNil(t, found)
		require.Equal(t, inputURL.OriginalURL, found.OriginalURL)
		require.Equal(t, inputURL.ShortURL, found.ShortURL)
		require.True(t, found.ID > 0)
	})
	t.Run("Get NotFound", func(t *testing.T) {
		found, err := repo.GetByShortURL(ctx, "not_exists")
		require.Nil(t, found)
		require.ErrorIs(t, err, exceptions.ErrNotFound)
	})
}
