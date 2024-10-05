package main

import (
	"log"

	"github.com/RajNykDhulapkar/gotiny/handler"
	"github.com/RajNykDhulapkar/gotiny/shortener"
	"github.com/RajNykDhulapkar/gotiny/store"
	"github.com/gin-gonic/gin"
)

func main() {
	storeService, err := store.InitializeStore("localhost:6379", "", 1)
	if err != nil {
		log.Fatalf("Failed to initialize the store: %v", err)
	}

	base58Encoder := shortener.NewBase58Encoder()
	shortenerService := shortener.NewShortener(base58Encoder)

	h := handler.NewHandler(storeService, shortenerService)

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.POST("/create-short-url", h.CreateShortURL)
	r.GET("/:shortUrl", h.HandleShortURLRedirect)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start the web server: %v", err)
	}
}
