package postgres

import (
// "context"
// "fmt"
// "tet/internals/models"
// "time"
)

// func LoadChatlist(userId int) []models.Chatlist {
// 	var ChatLists []models.Chatlist

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	query := "select receiver_id,last_msg,created_at from chatlist where sender_id = $1 and is_group = FALSE"
// 	rows, err := pool.Query(ctx,query, userId)
// 	if err != nil {
// 		fmt.Println("error while fetching chatlist from db - ", err)
// 		return nil
// 	}
// 	for rows.Next() {
// 		var chat_list models.Chatlist
// 		rows.Scan(
// 			&chat_list.ContactId,
// 			&chat_list.LastMsg,
// 			&chat_list.CreatedAt,
// 		)
// 		_, chat_list.ContactName, _, _, _, _, err = FindUser(chat_list.ContactId)
// 		if err != nil {
// 			fmt.Print("error - ", err.Error())
// 			return nil
// 		}
// 		fmt.Println("chatlist - ", chat_list)
// 		ChatLists = append(ChatLists, chat_list)
// 	}
// 	return ChatLists
// }

// func AddLastMsgToChatlist(senderId int, receiverId int, content string, createdAt time.Time) {

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	query := "update chatlist set last_msg = $1 , created_at = $2  where ( sender_id = $3 and receiver_id =$4 ) or ( receiver_id = $3 and sender_id =$4 ) "
// 	_, err := pool.Exec(ctx,query, content, createdAt, senderId, receiverId)
// 	if err != nil {
// 		fmt.Println("error while updating messages to chatlist - ", err)
// 		return
// 	}
// }
