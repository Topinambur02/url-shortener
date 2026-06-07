package inmemory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/topinambur02/url-shortener/internal/model"
	"github.com/topinambur02/url-shortener/pkg/exceptions"
)

func TestInMemoryRepository_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("successful link creation", func(t *testing.T) {
		repo := NewInMemoryRepository()
		expected := model.URL{
			ShortURL:    "shrt",
			OriginalURL: "https://example.com/very-long-url",
		}

		result, err := repo.Create(ctx, expected)

		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, expected, *result)

		repo.mu.RLock()
		savedUrl, exists := repo.data[expected.ShortURL]
		repo.mu.RUnlock()

		require.True(t, exists, "the link was not saved to the repository data map")
		require.Equal(t, expected, savedUrl)
	})

	t.Run("conflict error when re-creating", func(t *testing.T) {
		repo := NewInMemoryRepository()
		url := model.URL{
			ShortURL:    "duplicate",
			OriginalURL: "https://example.com",
		}

		_, err := repo.Create(ctx, url)
		require.NoError(t, err, "setup failed: expected no error on first creation")

		_, err = repo.Create(ctx, url)

		require.ErrorIs(t, err, exceptions.ErrConflict)
	})
}

func TestInMemoryRepository_GetByShortUrl(t *testing.T) {
	ctx := context.Background()

	t.Run("successful retrieval of an existing link", func(t *testing.T) {
		repo := NewInMemoryRepository()
		existing := model.URL{
			ShortURL:    "findme",
			OriginalURL: "https://find-this-site.com",
		}

		repo.mu.Lock()
		repo.data[existing.ShortURL] = existing
		repo.mu.Unlock()

		result, err := repo.GetByShortURL(ctx, existing.ShortURL)

		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, existing, *result)
	})

	t.Run("ErrNotFound error if there is no link", func(t *testing.T) {
		repo := NewInMemoryRepository()

		_, err := repo.GetByShortURL(ctx, "non-existent-key")

		require.ErrorIs(t, err, exceptions.ErrNotFound)
	})
}
