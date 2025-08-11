package postgres

import (
	// "database/sql"
	"context"
	"database/sql"
	"fmt"
	"tet/internals/models"
	"tet/internals/services"
	"tet/internals/storage/redis"
	"time"

	_ "github.com/lib/pq"
)

func StoreMessagesPostDB(message models.Message) int64 {
	redis_client := redis.GiveRedisConnection()
	var msgId int64
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	services.CheckForMessagePartition(redis_client, pool)
	if message.Type == "text/plain" {
		query := "insert into all_messages(sender_id,receiver_id,type,content,created_at) values($1,$2,$3,$4,$5) returning msg_id"
		err := pool.QueryRow(ctx, query, message.SenderID, message.ReceiverID, message.Type, message.Content, message.CreatedAt).Scan(&msgId)
		if err != nil {
			fmt.Println("error while inserting the messages - ", err)
			return 0
		}
		return msgId
	}
	return 0

}

func LoadChatMessagesPDb(userID int, contactID int, limit int, offset int) ([]models.Message, error) {
	var AllMessages []models.Message

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "select * from messages where sender_id =$1 and receiver_id = $2 order by msg_id desc limit $3 offset $4 "
	rows, err := pool.Query(ctx, query, userID, contactID, limit, offset)
	if err == sql.ErrNoRows {
		fmt.Println("no messages ")
		return nil, fmt.Errorf("Empty chat")
	}
	for rows.Next() {
		var message models.Message
		rows.Scan(
			&message.ID,
			&message.SenderID,
			&message.ReceiverID,
			&message.Type,
			&message.Content,
			&message.CreatedAt,
		)
		fmt.Println("message - ", message)
		AllMessages = append(AllMessages, message)
	}
	fmt.Printf("message for the sender_id - %v is %v", userID, AllMessages)
	return AllMessages, nil
}
