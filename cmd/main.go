package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/RajNykDhulapkar/gotiny/internals/cache"
	"github.com/RajNykDhulapkar/gotiny/internals/data"
	"github.com/RajNykDhulapkar/gotiny/internals/handler"
	"github.com/RajNykDhulapkar/gotiny/internals/rangeallocator"
	"github.com/RajNykDhulapkar/gotiny/internals/shortener"
	"github.com/gin-gonic/gin"
)

func main() {
	clientConfig := &rangeallocator.ClientConfig{
		Address:     os.Getenv("RANGE_ALLOCATOR_ADDRESS"),
		DialTimeout: 5 * time.Second,
	}

	client, err := rangeallocator.NewClient(clientConfig)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	mongodbConfig := &data.Config{
		URI:        os.Getenv("MONGODB_URI"),
		Database:   os.Getenv("MONGODB_DATABASE"),
		Collection: "urls",
	}

	repository, err := data.NewMongoRepository(mongodbConfig)
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}
	defer repository.Close(context.Background())

	managerConfig := &rangeallocator.RangeManagerConfig{
		ServiceID: os.Getenv("SERVICE_ID"),
		RangeSize: 1000,
		Region:    "default",
	}

	manager := rangeallocator.NewRangeManager(client, managerConfig)

	cacheConfig := &cache.Config{
		Address:      os.Getenv("REDIS_URL"),
		Password:     os.Getenv("REDIS_PASSWORD"),
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

	shortenerConfig := &shortener.Config{
		ServiceID: os.Getenv("SERVICE_ID"),
		RangeSize: 1000,
	}
	shortenerService := shortener.NewShortener(manager, shortenerConfig)

	h := handler.NewHandler(urlCache, shortenerService, repository)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.POST("/create-short-url", h.CreateShortURL)
	r.GET("/:shortUrl", h.HandleShortURLRedirect)
	r.GET("/urls/:userId", h.GetURLsByUserID) // Add this new route

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start the web server: %v", err)
	}
}
