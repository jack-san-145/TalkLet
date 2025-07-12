package redis

import (
	"fmt"
	"strconv"
	"tet/internals/models"
	"time"

	"github.com/redis/go-redis/v9"
)

func FindSessionRdb(session_id string) (int, bool, error) {
	userId_str, err := rdb.HGet(ctx, session_id, "user_id").Result()
	if err == redis.Nil {
		fmt.Println("session not found")
		return 0, false, nil
	} else if err != nil {
		return 0, false, fmt.Errorf("error while finding session_id in redis - %v", err)

	} else {
		fmt.Println("session found", userId_str)
		userId_str_int, _ := strconv.Atoi(userId_str)
		return userId_str_int, true, nil
	}
}

func SetSessionToRdb(session models.Session) {
	err := rdb.HSet(ctx, session.Session_id, "user_id", session.User_id, "expires_at", session.Expires_at).Err()
	if err != nil {
		fmt.Println("error while inserting session to redis - ", err)
	}
	err = rdb.Expire(ctx, session.Session_id, session.Expires_at.Sub(time.Now())).Err()
	if err != nil {
		fmt.Println("error while setting session expiry - ", err)
	}
	user_id, _, err := FindSessionRdb(session.Session_id)
	if err != nil {
		fmt.Println("error in insert - ", err)
		return
	}
	fmt.Println("user_id - ", user_id)
}

func DeleteSessionRdb(session_id string) {
	err := rdb.Del(ctx, session_id).Err()
	if err != nil {
		fmt.Println("error while deleting session in redis - ", err)
	}
}
