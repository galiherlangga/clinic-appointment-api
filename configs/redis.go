package configs

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func ConnectRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", AppConfig.RedisHost, AppConfig.RedisPort),
		Password: AppConfig.RedisPass,
		DB:       0, // Use default DB
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v. Logout/Blacklisting might not work.", err)
	} else {
		log.Println("Redis connection established")
	}
}
