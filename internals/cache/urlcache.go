package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/RajNykDhulapkar/gotiny/pkg/interfaces"
)

var (
	ErrKeyNotFound = fmt.Errorf("key not found in cache")
	ErrConnection  = fmt.Errorf("cache connection error")
)

type urlCacheImpl struct {
	cache interfaces.CachePort
}

func NewUrlCache(cache interfaces.CachePort) interfaces.UrlCache {
	return &urlCacheImpl{
		cache: cache,
	}
}

func (u *urlCacheImpl) SaveUrlMapping(ctx context.Context, shortUrl, originalUrl string, duration time.Duration) error {
	return u.cache.Set(ctx, shortUrl, originalUrl, duration)
}

func (u *urlCacheImpl) GetOriginalUrl(ctx context.Context, shortUrl string) (string, error) {
	return u.cache.Get(ctx, shortUrl)
}
