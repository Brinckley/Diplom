package redis

import (
	"context"
	"encoding/json"
	"events-subscriber/internal/structs"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"time"
)

var cRedisPassword string
var cRedisPort string

var ctx context.Context

func InitRedis() {
	ctx = context.Background()
	cRedisPassword = os.Getenv("REDIS_PASSWORD")
	cRedisPort = os.Getenv("REDIS_PORT")
}

func GetSubscriber() (*redis.Client, error) {
	var redisClient = redis.NewClient(&redis.Options{
		Addr:     cRedisPort,
		Password: cRedisPassword,
	})

	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		// Sleep for 3 seconds and wait for Redis to initialize
		time.Sleep(3 * time.Second)
		err := redisClient.Ping(context.Background()).Err()
		if err != nil {
			return nil, err
		}
	}

	return redisClient, err
}

func GetSubscriptionData(client *redis.Client, key string) error {
	val, err := client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	var eventList structs.Events
	err = json.Unmarshal([]byte(val), &eventList)
	if err != nil {
		return err
	}

	fmt.Println("value : ", eventList)
	return nil
}
