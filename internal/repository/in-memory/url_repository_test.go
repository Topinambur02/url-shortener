package inmemory

import (
	"context"
	"errors"
	"testing"

	"github.com/topinambur02/url-shortener/internal/model"
	"github.com/topinambur02/url-shortener/pkg/exceptions"
)

func TestUrlRepository(t *testing.T) {
	t.Run("TestUrlRepositoryImpl_Create", func(t *testing.T) {
		ctx := context.Background()

		t.Run("successful link creation", func(t *testing.T) {
			repo := NewUrlRepository()
			expectedUrl := model.Url{
				ShortUrl:    "shrt",
				OriginalUrl: "https://example.com/very-long-url",
			}

			result, err := repo.Create(ctx, expectedUrl)

			if err != nil {
				t.Fatalf("Expected error nil, received: %v", err)
			}

			if result == nil {
				t.Fatal("expected a model object, got nil")
			}
			if result.ShortUrl != expectedUrl.ShortUrl || result.OriginalUrl != expectedUrl.OriginalUrl {
				t.Errorf("expected %v, received %v", expectedUrl, *result)
			}

			repo.mu.RLock()
			savedUrl, exists := repo.data[expectedUrl.ShortUrl]
			repo.mu.RUnlock()

			if !exists {
				t.Error("the link was not saved to the repository")
			}
			if savedUrl.OriginalUrl != expectedUrl.OriginalUrl {
				t.Errorf("An incorrect value was saved in the map: %s", savedUrl.OriginalUrl)
			}
		})

		t.Run("conflict error when re-creating", func(t *testing.T) {
			repo := NewUrlRepository()
			url := model.Url{
				ShortUrl:    "duplicate",
				OriginalUrl: "https://example.com",
			}

			_, err := repo.Create(ctx, url)
			if err != nil {
				t.Fatalf("error while preparing the test: %v", err)
			}

			_, err = repo.Create(ctx, url)

			if !errors.Is(err, exceptions.ErrConflict) {
				t.Errorf("Expected error %v, received: %v", exceptions.ErrConflict, err)
			}
		})
	})
	t.Run("TestUrlRepositoryImpl_GetByShortUrl", func(t *testing.T) {
		ctx := context.Background()

		t.Run("successful retrieval of an existing link", func(t *testing.T) {
			repo := NewUrlRepository()
			existingUrl := model.Url{
				ShortUrl:    "findme",
				OriginalUrl: "https://find-this-site.com",
			}

			repo.data[existingUrl.ShortUrl] = existingUrl

			result, err := repo.GetByShortUrl(ctx, existingUrl.ShortUrl)

			if err != nil {
				t.Fatalf("Expected error nil, received: %v", err)
			}
			if result == nil {
				t.Fatal("Expected a model object, got nil")
			}
			if result.ShortUrl != existingUrl.ShortUrl || result.OriginalUrl != existingUrl.OriginalUrl {
				t.Errorf("expected %v, received %v", existingUrl, *result)
			}
		})

		t.Run("ErrNotFound error if there is no link", func(t *testing.T) {
			repo := NewUrlRepository()

			_, err := repo.GetByShortUrl(ctx, "non-existent-key")

			if !errors.Is(err, exceptions.ErrNotFound) {
				t.Errorf("Expected error %v, received: %v", exceptions.ErrNotFound, err)
			}
		})
	})
}
