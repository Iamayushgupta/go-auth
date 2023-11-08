package util

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

var ctx = context.Background()

func SetSession(redisClient *redis.Client, username string, sessionToken string) error {
	key := sessionToken
	err := redisClient.Set(ctx, key, username, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to set session in Redis")
	}
	return nil
}

func GetUsernameFromSession(redisClient *redis.Client, sessionToken string) (string, error) {
	key := sessionToken
	username, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get username from session")
	}
	return username, nil
}

func DeleteSession(redisClient *redis.Client, sessionToken string) error {
	key := sessionToken
	err := redisClient.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete session")
	}
	return nil
}

func GenerateRandomString(length int) (string, error) {
	if length <= 0 {
		return "", nil
	}
	// Calculate the required byte size
	byteSize := (length * 6) / 8
	if (length*6)%8 > 0 {
		byteSize++
	}
	bytes := make([]byte, byteSize)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	// Encode the random bytes into a base64 string
	randomString := base64.URLEncoding.EncodeToString(bytes)
	return randomString[:length], nil
}
