package services

import (
	"time"

	"github.com/galiherlangga/clinic-appointment/configs"
)

type BlacklistService interface {
	BlacklistToken(token string, expiration time.Duration) error
	IsTokenBlacklisted(token string) (bool, error)
}

type blacklistService struct{}

func NewBlacklistService() BlacklistService {
	return &blacklistService{}
}

func (s *blacklistService) BlacklistToken(token string, expiration time.Duration) error {
	key := configs.AppConfig.RedisPrefix + "blacklist:" + token
	return configs.RedisClient.Set(configs.Ctx, key, "true", expiration).Err()
}

func (s *blacklistService) IsTokenBlacklisted(token string) (bool, error) {
	key := configs.AppConfig.RedisPrefix + "blacklist:" + token
	val, err := configs.RedisClient.Get(configs.Ctx, key).Result()
	if err != nil {
		return false, nil // Token not found in blacklist
	}
	return val == "true", nil
}
