package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb *redis.Client

func CreateRedisClient() {
	rdb = redis.NewClient(&redis.Options{
		Addr:            "localhost:6379",
		MaxIdleConns:    10,
		ConnMaxLifetime: time.Hour,
	})
	err := rdb.Ping(ctx).Err()
	if err != nil {
		fmt.Println("redis client connection faliure - ", err)
		return
	}
	fmt.Println("redis connection successfull", err)
}

func GiveRedisConnection() *redis.Client {
	return rdb
}
