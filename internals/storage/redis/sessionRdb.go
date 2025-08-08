package redis

import (
	"fmt"
	"tet/internals/models"
	"time"

	"github.com/redis/go-redis/v9"
)

func FindSessionRdb(session_id string) (string, bool, error) {
	roll_no, err := rdb.HGet(ctx, session_id, "roll_no").Result()
	if err == redis.Nil {
		fmt.Println("session not found")
		return "", false, nil
	} else if err != nil {
		return "", false, fmt.Errorf("error while finding session_id in redis - %v", err)

	} else {
		fmt.Println("session found", roll_no)
		return roll_no, true, nil
	}
}

func SetSessionToRdb(session models.Session) {
	err := rdb.HSet(ctx, session.Session_id, "roll_no", session.Roll_no, "expires_at", session.Expires_at).Err()
	if err != nil {
		fmt.Println("error while inserting session to redis - ", err)
	}
	err = rdb.Expire(ctx, session.Session_id, session.Expires_at.Sub(time.Now())).Err()
	if err != nil {
		fmt.Println("error while setting session expiry - ", err)
	}
	roll_no, _, err := FindSessionRdb(session.Session_id)
	if err != nil {
		fmt.Println("error in insert - ", err)
		return
	}
	fmt.Println("roll_no - ", roll_no)
}

func DeleteSessionRdb(session_id string) {
	err := rdb.Del(ctx, session_id).Err()
	if err != nil {
		fmt.Println("error while deleting session in redis - ", err)
	}
}
