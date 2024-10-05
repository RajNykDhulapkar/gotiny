package store

import (
	"context"
	"fmt"
	"time"

	"github.com/RajNykDhulapkar/gotiny/interfaces"
	"github.com/redis/go-redis/v9"
)

type StorageService struct {
	redisClient *redis.Client
}

const CacheDuration = 6 * time.Hour

func InitializeStore(addr, password string, db int) (interfaces.StorageServiceInterface, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("error initializing Redis: %v", err)
	}

	fmt.Println("Redis started successfully")
	return &StorageService{redisClient: redisClient}, nil
}

func (s *StorageService) SaveUrlMapping(ctx context.Context, shortUrl, originalUrl, userId string) error {
	err := s.redisClient.Set(ctx, shortUrl, originalUrl, CacheDuration).Err()
	if err != nil {
		return fmt.Errorf("error saving URL mapping: %v", err)
	}
	return nil
}

func (s *StorageService) RetrieveOriginalUrl(ctx context.Context, shortUrl string) (string, error) {
	originalUrl, err := s.redisClient.Get(ctx, shortUrl).Result()
	if err != nil {
		return "", fmt.Errorf("error retrieving original URL: %v", err)
	}
	return originalUrl, nil
}
