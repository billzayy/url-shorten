package db

import (
	"github.com/go-redis/redis"
)

// var ctx = context.Background()

func InitRedisClient() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // Redis server address
		Password: "",           // No password set
		DB:       0,            // Use default DB
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		return &redis.Client{}, err
	}

	return rdb, nil
}
