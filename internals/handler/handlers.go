package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/RajNykDhulapkar/gotiny/pkg/interfaces"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	cache      interfaces.UrlCache
	shortener  interfaces.ShortenerInterface
	repository interfaces.URLRepository
}

type UrlCreationRequest struct {
	LongUrl string `json:"long_url" binding:"required"`
	UserId  string `json:"user_id" binding:"required"`
}

type UrlsByUserResponse struct {
	Message string                  `json:"message"`
	Data    []*interfaces.URLEntity `json:"data"`
}

func NewHandler(cache interfaces.UrlCache, shortener interfaces.ShortenerInterface, repository interfaces.URLRepository) interfaces.HandlerInterface {
	return &Handler{
		cache:      cache,
		shortener:  shortener,
		repository: repository,
	}
}

func (h *Handler) CreateShortURL(c *gin.Context) {
	var creationRequest UrlCreationRequest
	if err := c.ShouldBindJSON(&creationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortUrl, err := h.shortener.GenerateShortLink(c, creationRequest.LongUrl, creationRequest.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	urlEntity := &interfaces.URLEntity{
		ShortURL:    shortUrl,
		OriginalURL: creationRequest.LongUrl,
		UserID:      creationRequest.UserId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ClickCount:  0,
	}

	if err := h.repository.Save(c, urlEntity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save URL"})
		return
	}

	if err := h.cache.SaveUrlMapping(c, shortUrl, creationRequest.LongUrl, 6*time.Hour); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	host := "gotiny.fun/"
	c.JSON(200, gin.H{
		"message":   "short url created successfully",
		"short_url": host + shortUrl,
	})
}

func (h *Handler) HandleShortURLRedirect(c *gin.Context) {
	shortUrl := c.Param("shortUrl")

	originalUrl, err := h.cache.GetOriginalUrl(c, shortUrl)

	if err == nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = h.repository.IncrementClickCount(ctx, shortUrl)
		}()

		c.JSON(http.StatusOK, gin.H{
			"message":      "URL retrieved successfully",
			"original_url": originalUrl,
		})
		return
	}

	urlEntity, err := h.repository.FindByShortURL(c, shortUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve URL"})
		return
	}
	if urlEntity == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	_ = h.cache.SaveUrlMapping(c, shortUrl, urlEntity.OriginalURL, 6*time.Hour)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = h.repository.IncrementClickCount(ctx, shortUrl)
	}()

	c.JSON(http.StatusOK, gin.H{
		"message":      "URL retrieved successfully",
		"original_url": urlEntity.OriginalURL,
	})
}

func (h *Handler) GetURLsByUserID(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
		})
		return
	}

	urls, err := h.repository.FindByUserID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch URLs",
		})
		return
	}

	if len(urls) == 0 {
		c.JSON(http.StatusOK, UrlsByUserResponse{
			Message: "No URLs found for this user",
			Data:    []*interfaces.URLEntity{},
		})
		return
	}

	c.JSON(http.StatusOK, UrlsByUserResponse{
		Message: "URLs retrieved successfully",
		Data:    urls,
	})
}
