package repository

import (
	"context"
	"fmt"

	"github.com/topinambur02/url-shortener/internal/config"
)

type Factory func(ctx context.Context, cfg *config.Config) (StoreRepository, error)

var registry = make(map[string]Factory)

func Register(name string, f Factory) {
	registry[name] = f
}

func InitRepository(ctx context.Context, storageType string, cfg *config.Config) (StoreRepository, error) {
	f, ok := registry[storageType]
	if !ok {
		return nil, fmt.Errorf("unknown storage type: %s", storageType)
	}
	return f(ctx, cfg)
}