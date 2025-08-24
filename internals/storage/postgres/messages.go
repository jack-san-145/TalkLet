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

func LoadOTOChatMessagesPDb(userID string, contactID string, limit int, offset int) ([]models.Message, error) {
	var AllMessages []models.Message
	var (
		meta_data_validate any
		meta_data          models.MetaData
		created_at_time    time.Time
	)

	services.Find_dept_from_rollNo(userID)
	query := "select * from all_messages where (sender_id =$1 and receiver_id = $2) or (sender_id =$3 and receiver_id = $4) order by msg_id desc limit $5 offset $6 "
	rows, err := pool.Query(context.Background(), query, userID, contactID, contactID, userID, limit, offset)
	if err == sql.ErrNoRows {
		fmt.Println("no messages")
		return nil, fmt.Errorf("empty chat")
	}

	// type Message struct {
	//     ID         int64  `json:"msg_id"`
	//     SenderID   string `json:"sender_id"`
	//     ReceiverID string `json:"receiver_id"`
	//     Type       string `json:"type"`
	//     Content    string `json:"content"`
	//     CreatedAt  string `json:"created_at"`
	//     IsAck      string `json:"is_ack"`
	//     Status     string `json:"status"`
	// }

	for rows.Next() {
		var message models.Message
		err := rows.Scan(
			&message.ID,
			&message.SenderID,
			&message.ReceiverID,
			&message.Type,
			&message.Content,
			&meta_data_validate,
			&created_at_time,
			&message.Status,
		)
		message.CreatedAt = created_at_time.Format("2006-01-02 15:04:05")
		if meta_data_validate != nil {
			value, ok := meta_data_validate.(models.MetaData)
			if ok {
				meta_data = value
			}
			fmt.Println("meta data after type assertion - ", meta_data)
			message.MetaData = meta_data
		}
		if err != nil {
			fmt.Println("error while fetching the message history - ", err)
		}
		fmt.Println("meta_data - ", meta_data)
		fmt.Println("message - ", message)
		AllMessages = append(AllMessages, message)
	}
	fmt.Printf("message for the sender_id - %v is %v", userID, AllMessages)
	return AllMessages, nil
}
