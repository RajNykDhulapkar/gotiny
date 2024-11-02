package handler

import (
	"net/http"

	"github.com/RajNykDhulapkar/gotiny/pkg/interfaces"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Store     interfaces.StorageServiceInterface
	Shortener interfaces.ShortenerInterface
}

type UrlCreationRequest struct {
	LongUrl string `json:"long_url" binding:"required"`
	UserId  string `json:"user_id" binding:"required"`
}

func NewHandler(store interfaces.StorageServiceInterface, shortener interfaces.ShortenerInterface) interfaces.HandlerInterface {
	return &Handler{
		Store:     store,
		Shortener: shortener,
	}
}

func (h *Handler) CreateShortURL(c *gin.Context) {
	var creationRequest UrlCreationRequest
	if err := c.ShouldBindJSON(&creationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortUrl, err := h.Shortener.GenerateShortLink(creationRequest.LongUrl, creationRequest.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.Store.SaveUrlMapping(c, shortUrl, creationRequest.LongUrl, creationRequest.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	host := "http://localhost:8080/"
	c.JSON(200, gin.H{
		"message":   "short url created successfully",
		"short_url": host + shortUrl,
	})
}

func (h *Handler) HandleShortURLRedirect(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	initialUrl, err := h.Store.RetrieveOriginalUrl(c, shortUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(302, initialUrl)
}
