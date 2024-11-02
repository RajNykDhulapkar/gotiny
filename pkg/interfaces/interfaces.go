package interfaces

import (
	"context"
	"time"

	"github.com/RajNykDhulapkar/gotiny-range-allocator/pkg/pb"
	"github.com/gin-gonic/gin"
)

type HandlerInterface interface {
	CreateShortURL(c *gin.Context)
	HandleShortURLRedirect(c *gin.Context)
}

type ShortenerInterface interface {
	GenerateShortLink(initialLink string, userId string) (string, error)
}

type Base58EncoderInterface interface {
	Encode(input []byte) (string, error)
}

type CachePort interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Close() error
}

type UrlCache interface {
	SaveUrlMapping(ctx context.Context, shortUrl, originalUrl string, duration time.Duration) error
	GetOriginalUrl(ctx context.Context, shortUrl string) (string, error)
}

type RangeAllocatorPort interface {
	pb.RangeAllocatorServer
}
