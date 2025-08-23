package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func Set_OTP_to_redis(email string, otp string) {

	err := rdb.Set(context.Background(), email, otp, time.Minute).Err()
	if err != nil {
		fmt.Println("error while adding the otp to the redis - ", err)
		return
	}
}

func Get_OTP_from_redis(email string) (string, error) {
	otp, err := rdb.Get(context.Background(), email).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("returned nil")
	} else if err != nil {
		return "", fmt.Errorf("Invalid OTP - %v ", err)
	} else {
		return otp, nil
	}
}
