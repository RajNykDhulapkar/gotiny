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

type Base62EncoderPort interface {
	Encode(num int64) string
	Decode(encoded string) (int64, error)
}

type RangeAllocatorPort interface {
	GetNextID(ctx context.Context) (int64, error)
	GetCurrentRange() *pb.Range
}

type URLEntity struct {
	ID          string    `bson:"_id,omitempty"`
	ShortURL    string    `bson:"short_url"`
	OriginalURL string    `bson:"original_url"`
	UserID      string    `bson:"user_id"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
	ClickCount  int64     `bson:"click_count"`
}

type ShortenerInterface interface {
	GenerateShortLink(ctx context.Context, originalURL string, userID string) (string, error)
}

type URLRepository interface {
	Save(ctx context.Context, url *URLEntity) error
	FindByShortURL(ctx context.Context, shortURL string) (*URLEntity, error)
	FindByOriginalURL(ctx context.Context, originalURL string) (*URLEntity, error)
	IncrementClickCount(ctx context.Context, shortURL string) error
	FindByUserID(ctx context.Context, userID string) ([]*URLEntity, error)
	Close(ctx context.Context) error
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
