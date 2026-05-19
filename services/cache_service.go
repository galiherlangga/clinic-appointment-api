package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheService interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	DeletePattern(ctx context.Context, pattern string) error
}

type cacheService struct {
	client *redis.Client
}

func NewCacheService(client *redis.Client) CacheService {
	return &cacheService{client: client}
}

func (s *cacheService) Get(ctx context.Context, key string, dest interface{}) error {
	if s.client == nil {
		return redis.Nil
	}
	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func (s *cacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if s.client == nil {
		return nil
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, key, data, expiration).Err()
}

func (s *cacheService) Delete(ctx context.Context, key string) error {
	if s.client == nil {
		return nil
	}
	return s.client.Del(ctx, key).Err()
}

func (s *cacheService) DeletePattern(ctx context.Context, pattern string) error {
	if s.client == nil {
		return nil
	}
	iter := s.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := s.client.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}
