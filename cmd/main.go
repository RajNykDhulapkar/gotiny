package main

import (
	"context"
	"log"
	"time"

	"github.com/RajNykDhulapkar/gotiny/internals/cache"
	"github.com/RajNykDhulapkar/gotiny/internals/handler"
	"github.com/RajNykDhulapkar/gotiny/internals/rangeallocator"
	"github.com/RajNykDhulapkar/gotiny/internals/shortener"
	"github.com/gin-gonic/gin"
)

func main() {
	clientConfig := &rangeallocator.ClientConfig{
		Address:     "localhost:50051",
		DialTimeout: 5 * time.Second,
	}

	client, err := rangeallocator.NewClient(clientConfig)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	managerConfig := &rangeallocator.RangeManagerConfig{
		ServiceID: "url-shortener",
		RangeSize: 1000,
		Region:    "default",
	}

	manager := rangeallocator.NewRangeManager(client, managerConfig)

	id, err := manager.GetNextID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get next ID: %v", err)
	}

	log.Printf("Got next ID: %d", id)

	cacheConfig := &cache.Config{
		Address:      "localhost:6379",
		Password:     "",
		DB:           0,
		DialTimeout:  5 * time.Second,
		WriteTimeout: 2 * time.Second,
		ReadTimeout:  2 * time.Second,
	}

	redisAdapter, err := cache.NewRedisAdapter(cacheConfig)
	if err != nil {
		log.Fatalf("Failed to create Redis adapter: %v", err)
	}
	defer redisAdapter.Close()

	urlCache := cache.NewUrlCache(redisAdapter)
	if err != nil {
		log.Fatalf("Failed to initialize the store: %v", err)
	}

	base58Encoder := shortener.NewBase58Encoder()
	shortenerService := shortener.NewShortener(base58Encoder)

	h := handler.NewHandler(urlCache, shortenerService)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.POST("/create-short-url", h.CreateShortURL)
	r.GET("/:shortUrl", h.HandleShortURLRedirect)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start the web server: %v", err)
	}
}
