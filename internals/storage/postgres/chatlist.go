package postgres

import (
	"context"
	"fmt"
	"tet/internals/models"
	"tet/internals/services"
	"time"
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
		var created_at_time time.Time
		var chat_list models.ChatlistToSend
		err := rows.Scan(
			&chat_list.ContactId,
			&chat_list.LastMsg,
			&created_at_time,
		)
		fmt.Println("created_at_time - ", created_at_time)

		chat_list.CreatedAt = created_at_time.Format("2006-01-02 15:04:05")
		if err != nil {
			fmt.Println("error while scanning the chatlist - ", err)
		}
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

func AddLastMsgToChatlist(senderId string, receiverId string, last_msg_id int64, content string, createdAt string) {

	sender_dept := services.FindDeptChatlistByRollno(senderId)
	receiver_dept := services.FindDeptChatlistByRollno(receiverId)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//to add last msg to the sender
	query := fmt.Sprintf(`update %v set last_msg_id = $1 ,last_msg = $2 , created_at = $3  where sender_id = $4 and receiver_id =$5 `, sender_dept)
	_, err := pool.Exec(ctx, query, last_msg_id, content, createdAt, senderId, receiverId)
	if err != nil {
		fmt.Println("error while updating messages to sender's chatlist - ", err)
		return
	}

	//to add last msg to the receiver
	query = fmt.Sprintf(`update %v set last_msg_id = $1,last_msg = $2 , created_at = $3  where sender_id = $4 and receiver_id =$5 `, receiver_dept)
	_, err = pool.Exec(ctx, query, last_msg_id, content, createdAt, receiverId, senderId)
	if err != nil {
		fmt.Println("error while updating messages to receiver's chatlist - ", err)
		return
	}
}

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

	query = fmt.Sprintf(`insert into %s(sender_id,receiver_id,is_group,group_id,last_msg,last_msg_id,first_msg_id) values($1,$2,$3,$4,$5,$6,$7)`, dept_table)
	pool.Exec(context.Background(), query,
		newContact.ContactId,
		newContact.UserID,
		newContact.IsGroup,
		newContact.GroupId,
		newContact.LastMsg,
		newContact.LastMsgId,
		newContact.FirstMsgId)
}
