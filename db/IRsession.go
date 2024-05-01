package db

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func ReadIRSession(id string) (string, bool, error) {
	ctx := context.Background()
	val, err := redisClient.Get(ctx, fmt.Sprintf("SESSION:%s", id)).Result()
	if err != nil {
		if err == redis.Nil {
			return "", false, nil
		} else {
			return "", false, err
		}
	}
	return val, true, nil
}

func WriteIRSession(id, token string) error {
	ctx := context.Background()
	return redisClient.Set(ctx, fmt.Sprintf("SESSION:%s", id), token, 6*time.Hour).Err()
}

func RefreshIRSession(id string) error {
	ctx := context.Background()
	return redisClient.Expire(ctx, fmt.Sprintf("SESSION:%s", id), 6*time.Hour).Err()
}
