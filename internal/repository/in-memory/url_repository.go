package inmemory

import (
	"context"
	"sync"

	"github.com/topinambur02/url-shortener/internal/model"
	"github.com/topinambur02/url-shortener/pkg/exceptions"
)

type UrlRepositoryImpl struct {
	mu   sync.RWMutex
	data map[string]model.Url
}

func NewUrlRepository() *UrlRepositoryImpl {
	return &UrlRepositoryImpl{
		data: make(map[string]model.Url),
	}
}

func (r *UrlRepositoryImpl) GetByShortUrl(ctx context.Context, shortUrl string) (*model.Url, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	url, exists := r.data[shortUrl]
	if !exists {
        return nil, exceptions.ErrNotFound
    }

	return &url, nil
}

func (r *UrlRepositoryImpl) Create(ctx context.Context, url model.Url) (*model.Url, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[url.ShortUrl]; exists {
        return nil, exceptions.ErrConflict
    }

	r.data[url.ShortUrl] = url

	return &url, nil
}
