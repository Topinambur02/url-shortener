package inmemory

import (
	"context"
	"sync"

	"github.com/topinambur02/url-shortener/internal/config"
	"github.com/topinambur02/url-shortener/internal/model"
	"github.com/topinambur02/url-shortener/internal/repository"
	"github.com/topinambur02/url-shortener/pkg/exceptions"
)

func init() {
	repository.Register("inmemory", func(ctx context.Context, cfg *config.Config) (repository.StoreRepository, error) {
		return NewInMemoryRepository(), nil
	})
}

type InMemoryRepository struct {
	mu   sync.RWMutex
	data map[string]model.URL
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		data: make(map[string]model.URL),
	}
}

func (r *InMemoryRepository) GetByShortURL(ctx context.Context, shortURL string) (*model.URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	url, exists := r.data[shortURL]

	if !exists {
		return nil, exceptions.ErrNotFound
	}

	return &url, nil
}

func (r *InMemoryRepository) Create(ctx context.Context, url model.URL) (*model.URL, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[url.ShortURL]; exists {
		return nil, exceptions.ErrConflict
	}

	r.data[url.ShortURL] = url

	return &url, nil
}
