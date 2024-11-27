package data

import (
	"context"
	"fmt"
	"time"

	"github.com/RajNykDhulapkar/gotiny/pkg/interfaces"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	URI        string
	Database   string
	Collection string
}

type mongoRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoRepository(cfg *Config) (interfaces.URLRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	collection := client.Database(cfg.Database).Collection(cfg.Collection)
	return &mongoRepository{
		client:     client,
		collection: collection,
	}, nil
}

func (r *mongoRepository) Save(ctx context.Context, url *interfaces.URLEntity) error {
	_, err := r.collection.InsertOne(ctx, url)
	if err != nil {
		return fmt.Errorf("failed to save URL: %w", err)
	}
	return nil
}

func (r *mongoRepository) FindByShortURL(ctx context.Context, shortURL string) (*interfaces.URLEntity, error) {
	var result interfaces.URLEntity
	err := r.collection.FindOne(ctx, bson.M{"short_url": shortURL}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find URL by short URL: %w", err)
	}
	return &result, nil
}

func (r *mongoRepository) FindByOriginalURL(ctx context.Context, originalURL string) (*interfaces.URLEntity, error) {
	var result interfaces.URLEntity
	err := r.collection.FindOne(ctx, bson.M{"original_url": originalURL}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find URL by original URL: %w", err)
	}
	return &result, nil
}

func (r *mongoRepository) IncrementClickCount(ctx context.Context, shortURL string) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"short_url": shortURL},
		bson.M{"$inc": bson.M{"click_count": 1}},
	)
	if err != nil {
		return fmt.Errorf("failed to increment click count: %w", err)
	}
	return nil
}

func (r *mongoRepository) FindByUserID(ctx context.Context, userID string) ([]*interfaces.URLEntity, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("failed to find URLs by user ID: %w", err)
	}
	defer cursor.Close(ctx)

	var results []*interfaces.URLEntity
	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode URLs: %w", err)
	}
	return results, nil
}

func (r *mongoRepository) Close(ctx context.Context) error {
	return r.client.Disconnect(ctx)
}
