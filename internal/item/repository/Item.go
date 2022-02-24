package repository

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sedat/redis-api/internal/item/model"
)

type Repository interface {
	Get(ctx context.Context, key string) (string, error)
	GetAllKeys(ctx context.Context) ([]string, error)
	GetAllItems(ctx context.Context) ([]model.Item, error)
	Set(ctx context.Context, item model.Item) error
	Flush(ctx context.Context) bool
}

type repository struct {
	client *redis.Client
}

func NewRepository(client *redis.Client) Repository {
	return repository{client}
}

func (r repository) Get(ctx context.Context, key string) (string, error) {
	date := time.Now().Format(time.RFC3339)
	value, err := r.client.Get(ctx, key).Result()
	if err != nil {
		log.Printf("Couldn't find key %s at time %s %s", key, date, err)
		return "", err
	}
	log.Printf("Found key %s at time %s", key, date)
	return value, nil
}

func (r repository) GetAllKeys(ctx context.Context) ([]string, error) {
	date := time.Now().Format(time.RFC3339)
	var keys []string
	log.Printf("Retrieving all keys from redis %s", date)
	iter := r.client.Scan(ctx, 0, "*", 0).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		log.Printf("Error retrieving all keys from redis %s %s", date, err)
		return []string{}, err
	}
	return keys, nil
}

func (r repository) GetAllItems(ctx context.Context) ([]model.Item, error) {
	date := time.Now().Format(time.RFC3339)
	var items []model.Item
	log.Printf("Getting all items from redis %s", date)
	keys, err := r.GetAllKeys(ctx)
	if err != nil {
		log.Printf("Error retrieving all keys from redis %s %s", date, err)
		return []model.Item{}, err
	}
	for _, key := range keys {
		value, err := r.Get(ctx, key)
		if err != nil {
			log.Printf("Error getting %s from redis %s", key, err)
			return []model.Item{}, err
		}
		items = append(items, model.Item{Key: key, Value: value})
	}
	return items, nil
}

func (r repository) Set(ctx context.Context, item model.Item) error {
	date := time.Now().Format(time.RFC3339)
	err := r.client.Set(ctx, item.Key, item.Value, 0).Err()
	if err != nil {
		log.Printf("Couldn't set key  %s at time %s %s", item.Key, date, err)
		return err
	}
	log.Printf("Succesfully set key: %s value: %s t time %s", item.Key, item.Value, date)
	return nil
}

func (r repository) Flush(ctx context.Context) bool {
	date := time.Now().Format(time.RFC3339)
	err := r.client.FlushAllAsync(ctx).Err()
	if err != nil {
		log.Printf("Couldn't flush the database %s error: %s", date, err)
		return false
	}
	log.Printf("Succesfully flushed database %s", date)
	return true
}
