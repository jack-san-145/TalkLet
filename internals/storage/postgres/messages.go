package postgres

import (
	// "database/sql"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"tet/internals/models"
	"tet/internals/services"
	"tet/internals/storage/minio"
	"tet/internals/storage/redis"

	"time"

	_ "github.com/lib/pq"
)

func Store_Privatechat_MessagesPostDB(message models.Message) int64 {
	redis_client := redis.GiveRedisConnection()
	var msgId int64
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	services.Check_Private_MessagePartition(redis_client, pool)

	if message.Type == "text/plain" {
		query := "insert into all_private_messages(sender_id,receiver_id,type,content,created_at) values($1,$2,$3,$4,$5) returning msg_id"
		err := pool.QueryRow(ctx, query, message.SenderID, message.ReceiverID, message.Type, message.Content, message.CreatedAt).Scan(&msgId)
		if err != nil {
			fmt.Println("error while inserting the messages - ", err)
			return 0
		}
		return msgId

	} else if message.Type == "file" || message.Type == "media" {

		// meta_data := fmt.Sprintf(`'{file_name : %s , file_size : %d , file_url : %s , mime_type : %s}'`, message.MetaData.FileName, message.MetaData.FileSize, message.MetaData.FileURL, message.MetaData.MimeType)
		meta_data := map[string]any{
			"file_name": message.MetaData.FileName,
			"file_size": message.MetaData.FileSize,
			"file_url":  message.MetaData.FileURL,
			"mime_type": message.MetaData.MimeType,
		}
		meta_data_json, err := json.Marshal(meta_data)
		if err != nil {
			fmt.Println("error while marshal meta data to json - ", err)
			return 0
		}

		query := "insert into all_private_messages(sender_id,receiver_id,type,content,meta_data,created_at) values($1,$2,$3,$4, $5::jsonb ,$6) returning msg_id"
		err = pool.QueryRow(ctx, query, message.SenderID, message.ReceiverID, message.Type, message.Content, string(meta_data_json), message.CreatedAt).Scan(&msgId)
		if err != nil {
			fmt.Println("error while inserting the messages - ", err)
			return 0
		}
		return msgId
	}
	return 0

}

func Store_Groupchat_MessagesPostDB(message models.Message) int64 {
	redis_client := redis.GiveRedisConnection()
	var msgId int64
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dept := services.Find_dept_from_groupId(message.GroupId) // finding the dept
	services.Check_Group_MessagePartition(redis_client, pool, dept)
	table_name := dept + "group_all_messages" //table name varies depends upon departments
	if message.Type == "text/plain" {
		query := fmt.Sprintf(`insert into %s(sender_id,group_id,type,content,created_at) values($1,$2,$3,$4,$5) returning msg_id`, table_name)
		err := pool.QueryRow(ctx, query, message.SenderID, message.GroupId, message.Type, message.Content, message.CreatedAt).Scan(&msgId)
		if err != nil {
			fmt.Println("error while inserting the messages to the group chat - ", err)
			return 0
		}
		return msgId

	} else if message.Type == "file" || message.Type == "media" {

		// meta_data := fmt.Sprintf(`'{file_name : %s , file_size : %d , file_url : %s , mime_type : %s}'`, message.MetaData.FileName, message.MetaData.FileSize, message.MetaData.FileURL, message.MetaData.MimeType)
		meta_data := map[string]any{
			"file_name": message.MetaData.FileName,
			"file_size": message.MetaData.FileSize,
			"file_url":  message.MetaData.FileURL,
			"mime_type": message.MetaData.MimeType,
		}
		meta_data_json, err := json.Marshal(meta_data)
		if err != nil {
			fmt.Println("error while marshal meta data to json - ", err)
			return 0
		}

		query := fmt.Sprintf(`insert into %s(sender_id,group_id,type,content,meta_data,created_at) values($1,$2,$3,$4, $5::jsonb ,$6) returning msg_id`, table_name)
		err = pool.QueryRow(ctx, query, message.SenderID, message.GroupId, message.Type, message.Content, string(meta_data_json), message.CreatedAt).Scan(&msgId)
		if err != nil {
			fmt.Println("error while inserting the messages to the group chat - ", err)
			return 0
		}
		return msgId
	}
	return 0

}

func Load_PrivateChatMessages_PDb(userID string, contactID string, limit int, offset int) ([]models.Message, error) {
	var AllMessages []models.Message
	var (
		// meta_data_validate any
		// meta_data          models.MetaData
		created_at_time time.Time
	)

	services.Find_dept_from_rollNo(userID)
	query := `select msg_id,receiver_id,type,content,
				coalesce(meta_data ->> 'file_name','') as file_name,
				coalesce((meta_data ->> 'file_size')::bigint,0) as file_size,
				coalesce(meta_data ->> 'mime_type','') as mime_type,
				created_at,
				status from all_private_messages where (sender_id =$1 and receiver_id = $2) or (sender_id =$3 and receiver_id = $4) order by msg_id desc limit $5 offset $6 `
	rows, err := pool.Query(context.Background(), query, userID, contactID, contactID, userID, limit, offset)
	if err == sql.ErrNoRows {
		fmt.Println("no messages")
		return nil, fmt.Errorf("empty chat")
	} else if err != nil {
		fmt.Println("error while fetching the message history from db - ", err)
	}
	fmt.Println("comming outside the rows.Next() ")
	for rows.Next() {
		fmt.Println("comming inside the rows.Next() ")
		var message models.Message
		message.SenderID = userID
		err := rows.Scan(
			&message.ID,
			&message.ReceiverID,
			&message.Type,
			&message.Content,
			&message.MetaData.FileName,
			&message.MetaData.FileSize,
			&message.MetaData.MimeType,
			// &message.MetaData.FileURL,
			// &meta_data_validate,
			&created_at_time,
			&message.Status,
		)
		message.CreatedAt = created_at_time.Format("2006-01-02 15:04:05")
		if message.MetaData.FileName != "" {
			minio.GetFile_private_chats(&message)
		}
		// if meta_data_validate != nil {
		// 	value, ok := meta_data_validate.(models.MetaData)
		// 	if ok {
		// 		meta_data = value
		// 	}
		// 	fmt.Println("meta data after type assertion - ", meta_data)
		// 	message.MetaData = meta_data
		// }
		if err != nil {
			fmt.Println("error while scanning the message history - ", err)
		}
		// fmt.Println("meta_data - ", meta_data)
		fmt.Println("message - ", message)
		AllMessages = append(AllMessages, message)
	}
	fmt.Printf("message for the sender_id - %v is %v", userID, AllMessages)
	return AllMessages, nil
}

func Load_GroupChatMessages_PDb(userID string, groupID string, limit int, offset int) ([]models.Message, error) {
	var AllMessages []models.Message
	var (
		dept string
		// meta_data_validate any
		// meta_data          models.MetaData
		created_at_time time.Time
	)
	student_or_staff := services.Find_staff_or_student_by_id(userID)
	if student_or_staff == "STUDENT" {
		dept, _, _ = services.Find_dept_from_rollNo(userID)
	} else if student_or_staff == "STAFF" {
		var err error
		dept, _, _, err = Find_dept_from_staff_id(userID)
		if err != nil {
			return []models.Message{}, err //empty slice and err
		}
	}

	dept_table := dept + "_group_all_messages"

	query := fmt.Sprintf(`select msg_id,group_id,type,content,
				coalesce(meta_data ->> 'file_name','') as file_name,
				coalesce((meta_data ->> 'file_size')::bigint,0) as file_size,
				coalesce(meta_data ->> 'mime_type','') as mime_type,
				created_at,
				status from %s where group_id = $1 order by msg_id desc limit $2 offset $3 `, dept_table)
	rows, err := pool.Query(context.Background(), query, groupID, limit, offset)
	if err == sql.ErrNoRows {
		fmt.Println("no group messages yet ")
		return nil, fmt.Errorf("empty group chat")
	} else if err != nil {
		fmt.Println("error while fetching the group chat history from db - ", err)
	}
	fmt.Println("comming outside the rows.Next() ")
	for rows.Next() {
		fmt.Println("comming inside the rows.Next() ")
		var message models.Message
		message.SenderID = userID
		err := rows.Scan(
			&message.ID,
			&message.GroupId,
			&message.Type,
			&message.Content,
			&message.MetaData.FileName,
			&message.MetaData.FileSize,
			&message.MetaData.MimeType,
			&created_at_time,
			&message.Status,
		)
		message.CreatedAt = created_at_time.Format("2006-01-02 15:04:05")
		if message.MetaData.FileName != "" {
			minio.GetFile_private_chats(&message)
		}
		if err != nil {
			fmt.Println("error while scanning the group chat history - ", err)
		}
		// fmt.Println("meta_data - ", meta_data)
		fmt.Println("message - ", message)
		AllMessages = append(AllMessages, message)
	}
	fmt.Printf("message for the group_id - %v is %v", groupID, AllMessages)
	return AllMessages, nil
}
