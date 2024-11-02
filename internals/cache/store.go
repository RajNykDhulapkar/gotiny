package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/RajNykDhulapkar/gotiny/pkg/interfaces"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Address      string        `yaml:"address"`
	Password     string        `yaml:"password"`
	DB           int           `yaml:"db"`
	DialTimeout  time.Duration `yaml:"dialTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`
}

func DefaultConfig() *Config {
	return &Config{
		Address:      "localhost:6379",
		Password:     "",
		DB:           0,
		DialTimeout:  5 * time.Second,
		WriteTimeout: 2 * time.Second,
		ReadTimeout:  2 * time.Second,
	}
}

type redisAdapter struct {
	client *redis.Client
}

func NewRedisAdapter(cfg *Config) (interfaces.CachePort, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Address,
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  cfg.DialTimeout,
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), cfg.DialTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &redisAdapter{
		client: client,
	}, nil
}

func (r *redisAdapter) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if err := r.client.Set(ctx, key, value, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set value in Redis: %w", err)
	}
	return nil
}

func (r *redisAdapter) Get(ctx context.Context, key string) (string, error) {
	value, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key not found: %s", key)
	}
	if err != nil {
		return "", fmt.Errorf("failed to get value from Redis: %w", err)
	}
	return value, nil
}

func (r *redisAdapter) Delete(ctx context.Context, key string) error {
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete key from Redis: %w", err)
	}
	return nil
}

func (r *redisAdapter) Close() error {
	return r.client.Close()
}
