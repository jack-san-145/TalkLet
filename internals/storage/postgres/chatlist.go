package postgres

import (
	"context"
	"fmt"
	"tet/internals/models"
	"tet/internals/services"
)

// "context"
// "fmt"
// "tet/internals/models"
// "time"

func LoadChatlist(userId string) []models.ChatlistToSend {
	var ChatLists []models.ChatlistToSend

	dept_table := services.FindDeptChatlistByRollno(userId)
	query := fmt.Sprintf(`select receiver_id,last_msg,created_at from %s where sender_id = $1 and is_group = FALSE`, dept_table)
	rows, err := pool.Query(context.Background(), query, userId)
	if err != nil {
		fmt.Println("error while fetching one to one chatlist from db - ", err)
		return nil
	}
	for rows.Next() {
		var chat_list models.ChatlistToSend
		rows.Scan(
			&chat_list.ContactId,
			&chat_list.LastMsg,
			&chat_list.CreatedAt,
		)
		chat_list.ContactName, _, _, err = FindUser(chat_list.ContactId)
		if err != nil {
			fmt.Print("error - ", err.Error())
			return nil
		}
		fmt.Println("chatlist - ", chat_list)
		ChatLists = append(ChatLists, chat_list)
	}
	return ChatLists
}

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

func AddTochatlist(newContact models.ChatlistForLocal, isGroup bool) {

	dept_table := services.FindDeptStudentByRollNo(newContact.UserID)
	query := fmt.Sprintf(`insert into %s(sender_id,receiver_id,is_group,group_id,last_msg,last_msg_id,first_msg_id) values($1,$2,$3,$4,$5,$6,$7)`, dept_table)
	pool.Exec(context.Background(), query,
		newContact.UserID,
		newContact.ContactId,
		newContact.IsGroup,
		newContact.GroupId,
		newContact.LastMsg,
		newContact.LastMsgId,
		newContact.FirstMsgId)
}
