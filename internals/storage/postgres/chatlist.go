package postgres

import (
	"context"
	"fmt"
	"tet/internals/models"
	"tet/internals/services"
	"time"
)

func LoadChatlist(userId string) []models.ChatlistToSend {
	var ChatLists []models.ChatlistToSend
	var (
		// dept           string
		// dept_table     string
		chatlist_table string
	)
	is_student_or_staff := services.Find_staff_or_student_by_id(userId)
	if is_student_or_staff == "STUDENT" {
		_, _, chatlist_table = services.Find_dept_from_rollNo(userId)
	} else if is_student_or_staff == "STAFF" {
		_, _, chatlist_table, _ = Find_dept_from_staff_id(userId)

	}

	query := fmt.Sprintf(`select receiver_id,last_msg,created_at,is_group,group_id,last_msg_id from %s where sender_id = $1`, chatlist_table)
	rows, err := pool.Query(context.Background(), query, userId)
	if err != nil {
		fmt.Println("error while fetching chatlist from db - ", err)
		return nil
	}

	for rows.Next() {
		var created_at_time time.Time
		var chat_list models.ChatlistToSend
		chat_list.UserId = userId
		err := rows.Scan(
			&chat_list.ContactId,
			&chat_list.LastMsg,
			&created_at_time,
			&chat_list.IsGroup,
			&chat_list.Group_id,
			&chat_list.LastMsgId)
		fmt.Println("created_at_time - ", created_at_time)

		chat_list.CreatedAt = created_at_time.Format("2006-01-02 15:04:05")
		if err != nil {
			fmt.Println("error while scanning the chatlist - ", err)
		}

		if chat_list.IsGroup {
			chat_list.GroupName = Find_groupname_by_groupid(chat_list.Group_id) //here contact name in the sense it as the group name
		} else {
			chat_list.ContactName, _, _, err = FindContact(chat_list.ContactId)
		}

		if err != nil {
			fmt.Print("error - ", err.Error())
			return nil
		}
		fmt.Println("chatlist - ", chat_list)
		ChatLists = append(ChatLists, chat_list)
	}
	return ChatLists
}

func AddLastMsgToChatlist_private_chat(senderId string, receiverId string, last_msg_id int64, content string, createdAt string) {

	sender_dept := services.FindDeptChatlistByRollno(senderId)
	receiver_dept := services.FindDeptChatlistByRollno(receiverId)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//to add last msg to the sender
	query := fmt.Sprintf(`update %s set last_msg_id = $1 ,last_msg = $2 , created_at = $3  where sender_id = $4 and receiver_id =$5 `, sender_dept)
	_, err := pool.Exec(ctx, query, last_msg_id, content, createdAt, senderId, receiverId)
	if err != nil {
		fmt.Println("error while updating messages to sender's chatlist - ", err)
		return
	}

	//to add last msg to the receiver
	query = fmt.Sprintf(`update %s set last_msg_id = $1,last_msg = $2 , created_at = $3  where sender_id = $4 and receiver_id =$5 `, receiver_dept)
	_, err = pool.Exec(ctx, query, last_msg_id, content, createdAt, receiverId, senderId)
	if err != nil {
		fmt.Println("error while updating messages to receiver's chatlist - ", err)
		return
	}
}

func AddLastMsgToChatlist_group_chat(message models.Message) {

	var exist bool
	dept_from_group := services.Find_dept_from_groupId(message.GroupId)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//to find the groupmembers
	All_group_members, err := Get_all_group_members(message.SenderID, dept_from_group)
	if err != nil {
		fmt.Println("error while getting group members - ", err)
		return
	}

	for _, group_member := range All_group_members {
		for group_member_id, dept_chatlist_table := range group_member {

			//checking the group is existing on m=everyone chatlist
			check_query := fmt.Sprintf(`select exists(select 1 from %s where (sender_id = $1 and group_id = $2) )`, dept_chatlist_table)
			err = pool.QueryRow(ctx, check_query, group_member_id, message.GroupId).Scan(&exist)
			if !exist {

				//if doesn't exists insert that group into chatlist
				insert_query := fmt.Sprintf(`insert into %s(sender_id,is_group,group_id,last_msg,last_msg_id,created_at) values($1,$2,$3,$4,$5,$6)`, dept_chatlist_table)
				_, err := pool.Exec(ctx, insert_query, group_member_id, true, message.GroupId, message.Content, message.GroupId, message.CreatedAt)
				if err != nil {
					fmt.Println("error while inserting data to chatlist - ", err)
					return
				}
			} else {

				//if group already exists then update last msg with new one
				query := fmt.Sprintf(`update %s set last_msg_id = $1 ,last_msg = $2 , created_at = $3  where (sender_id = $4 and group_id = $5)  `, dept_chatlist_table)
				_, err := pool.Exec(ctx, query, message.ID, message.Content, message.CreatedAt, group_member_id, message.GroupId)
				if err != nil {
					fmt.Println("error while updating messages to sender's chatlist - ", err)
					return
				}
			}

		}
	}

	//to add last msg to the sender
	// query := fmt.Sprintf(`update %s set last_msg_id = $1 ,last_msg = $2 , created_at = $3  where sender_id = $4 and receiver_id =$5 `, sender_dept)
	// _, err := pool.Exec(ctx, query, last_msg_id, content, createdAt, senderId, receiverId)
	// if err != nil {
	// 	fmt.Println("error while updating messages to sender's chatlist - ", err)
	// 	return
	// }

	// //to add last msg to the receiver
	// query = fmt.Sprintf(`update %s set last_msg_id = $1,last_msg = $2 , created_at = $3  where sender_id = $4 and receiver_id =$5 `, receiver_dept)
	// _, err = pool.Exec(ctx, query, last_msg_id, content, createdAt, receiverId, senderId)
	// if err != nil {
	// 	fmt.Println("error while updating messages to receiver's chatlist - ", err)
	// 	return
	// }
}

// func AddTochatlist(newContact models.ChatlistForLocal, isGroup bool) {

// 	// dept_table := services.FindDeptStudentByRollNo(newContact.UserID)
// 	query := fmt.Sprintf(`insert into %s(sender_id,receiver_id,is_group,group_id,last_msg,last_msg_id,first_msg_id) values($1,$2,$3,$4,$5,$6,$7)`, dept_table)
// 	pool.Exec(context.Background(), query,
// 		newContact.UserID,
// 		newContact.ContactId,
// 		newContact.IsGroup,
// 		newContact.GroupId,
// 		newContact.LastMsg,
// 		newContact.LastMsgId,
// 		newContact.FirstMsgId)

// 	query = fmt.Sprintf(`insert into %s(sender_id,receiver_id,is_group,group_id,last_msg,last_msg_id,first_msg_id) values($1,$2,$3,$4,$5,$6,$7)`, dept_table)
// 	pool.Exec(context.Background(), query,
// 		newContact.ContactId,
// 		newContact.UserID,
// 		newContact.IsGroup,
// 		newContact.GroupId,
// 		newContact.LastMsg,
// 		newContact.LastMsgId,
// 		newContact.FirstMsgId)
// }
