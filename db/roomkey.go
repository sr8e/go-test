package db

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/redis/go-redis/v9"
)

func SetRoomKey(roomName string) (string, error) {
	roomKey := generateRoomKey()
	_, exists, err := GetRoomName(roomKey)
	if err != nil {
		return "", err
	}

	for exists {
		roomKey = generateRoomKey()
		_, exists, err = GetRoomName(roomKey)
		if err != nil {
			return "", err
		}
	}

	// TODO: use watch
	ctx := context.Background()
	_, err = redisClient.TxPipelined(ctx, func(p redis.Pipeliner) error {
		p.Set(ctx, fmt.Sprintf("ROOMOF:%s", roomKey), roomName, 0)
		p.Set(ctx, fmt.Sprintf("ROOMKEYOF:%s", roomName), roomKey, 0)
		return nil
	})
	if err != nil {
		return "", err
	}
	return roomKey, nil
}

func GetRoomName(roomKey string) (string, bool, error) {
	ctx := context.Background()
	roomName, err := redisClient.Get(ctx, fmt.Sprintf("ROOMOF:%s", roomKey)).Result()
	if err != nil {
		if err == redis.Nil {
			return "", false, nil
		}
		return "", false, err
	}
	return roomName, true, nil
}

func ClearRoom(roomName string) error {
	ctx := context.Background()
	// TODO: use watch
	roomKey, err := redisClient.Get(ctx, fmt.Sprintf("ROOMKEYOF:%s", roomName)).Result()
	if err != nil {
		return err
	}
	_, err = redisClient.TxPipelined(ctx, func(p redis.Pipeliner) error {
		p.Del(ctx, fmt.Sprintf("ROOMOF:%s", roomKey))
		p.Del(ctx, fmt.Sprintf("ROOMKEYOF:%s", roomName))
		return nil
	})
	return err
}

func generateRoomKey() string {
	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	key := make([]rune, 4)

	for i := 0; i < 4; i++ {
		key[i] = letters[rand.Intn(len(letters))]
	}
	return string(key)
}
