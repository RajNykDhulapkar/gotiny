package interfaces

import (
	"context"

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

type StorageServiceInterface interface {
	SaveUrlMapping(ctx context.Context, shortUrl, originalUrl, userId string) error
	RetrieveOriginalUrl(ctx context.Context, shortUrl string) (string, error)
}
