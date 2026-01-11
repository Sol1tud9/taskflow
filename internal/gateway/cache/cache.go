package cache

import "context"

type Cache interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}) error
	Delete(ctx context.Context, key string) error
	Close() error
}

