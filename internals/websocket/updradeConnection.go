package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tet/internals/handlers"
	"tet/internals/models"
	"tet/internals/services"
	"tet/internals/storage/postgres"
	"time"

	"github.com/gorilla/websocket"
)

func UpgradeToWebsocket(w http.ResponseWriter, r *http.Request) {
	CookieFound, senderID := handlers.FindCookie(r)
	if !CookieFound {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	websocketConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("error while upgrade http to websocket - ", err)
		return
	}
	Set_ws_conn(senderID, websocketConn) //store the new ws connection to the ConnMap

	ConnMap.Range(func(key any, value any) bool {
		id, ok := key.(string)
		// conn:=value.(*websocket.Conn)
		fmt.Printf("id - %v is connected\n", id)
		return ok
	})
	// Get_ws_Conn(senderID)
	defer func() {
		fmt.Println("yeah deleting the connection")
		del_ws_conn(senderID) //to delete the ws_connection from the ConnMap
	}()

	listen_for_ws_msg(websocketConn, senderID)
}

func listen_for_ws_msg(websocketConn *websocket.Conn, senderID string) {

	var message_from_sender models.Message //this is my message which i want to send to my frd
	for {

		msg_type, msg, err := websocketConn.ReadMessage()
		msg_time := time.Now().Format("2006-01-02 15:04:05")
		message_from_sender.CreatedAt = msg_time
		message_from_sender.SenderID = senderID
		if err != nil {
			fmt.Println("error while reading ws- ", err)
			return
		}

		err = json.Unmarshal(msg, &message_from_sender)
		if err != nil {
			fmt.Println("error while unmarshal - ", err)
		}

		if message_from_sender.IsGroup {
			send_this_msg_to_group_chat(&message_from_sender, msg_type) //send ack for sender message
		} else {
			send_this_msg_to_private_chat(&message_from_sender, msg_type) //send ack for sender message
		}

	}
}

func send_this_msg_to_private_chat(message_from_sender *models.Message, msg_type int) {

	//sending ack to the sender itself

	// var message_to_receiver models.Message //this is the my message which is going to send to my frd
	message_from_sender.IsAck = "ack"

	// var temp models.Message
	if msg_type == 1 { // which is websocket.TextMessage(1)
		message_from_sender.Type = "text/plain"
		message_from_sender.ID = postgres.Store_Privatechat_MessagesPostDB(*message_from_sender)
		message_from_sender.Status = "sent"
	}

	go postgres.AddLastMsgToChatlist_private_chat(message_from_sender)
	msg_sent_by_sender_byte, err := json.Marshal(message_from_sender)
	if err != nil {
		fmt.Println("error while marshal ack - ", err)
		return
	}
	fmt.Printf(" ack - %+v\n ", message_from_sender)

	ws_conn, is_alive := Get_ws_Conn(message_from_sender.SenderID)
	if !is_alive {
		fmt.Println("websocket connection is not found")
		return
	}

	err = ws_conn.WriteMessage(msg_type, msg_sent_by_sender_byte)
	if err != nil {
		fmt.Println("error while sending ack to sender - ", err)
	}

	fmt.Println("received message  details after ack - ", message_from_sender)

	// sending the received msg

	// message_to_receiver = *message_from_sender

	// message_to_receiver.ID = message_from_sender.ID
	// message_to_receiver.SenderID = message_from_sender.SenderID
	// message_to_receiver.ReceiverID = message_from_sender.ReceiverID
	// message_to_receiver.Content = message_from_sender.Content
	// message_to_receiver.Type = message_from_sender.Type
	// message_to_receiver.CreatedAt = message_from_sender.CreatedAt
	sendMsgTo_receiver_private_chat(message_from_sender, msg_type) //send to frd

}

func sendMsgTo_receiver_private_chat(message_to_receiver *models.Message, msg_type int) {
	message_to_receiver.IsAck = "not-ack"
	sending_msg_byte, err := json.Marshal(message_to_receiver)
	if err != nil {
		fmt.Println("error while marshal sending_msg_byte - ", sending_msg_byte)
		return
	}
	fmt.Println("msg_send_to_frd - ", message_to_receiver)

	receiver_id, ok := Get_ws_Conn(message_to_receiver.ReceiverID)
	if !ok {
		fmt.Printf("%v goes to offline - ", message_to_receiver.ReceiverID)
		return
	}
	err = receiver_id.WriteMessage(msg_type, sending_msg_byte)
	if err != nil {
		fmt.Println("error while sending message_to_receiver private - ", err)
	}
}

// function to send ack for sender message to him itself
func send_this_msg_to_group_chat(message_from_sender *models.Message, msg_type int) {
	message_from_sender.IsAck = "ack"
	message_from_sender.SenderDept = services.Find_dept_from_groupId(message_from_sender.GroupId)
	// var temp models.Message
	if msg_type == 1 { // which is websocket.TextMessage(1)
		message_from_sender.Type = "text/plain"
		message_from_sender.ID = postgres.Store_Groupchat_MessagesPostDB(*message_from_sender)
		message_from_sender.Status = "sent"
	}

	go postgres.AddLastMsgToChatlist_group_chat(message_from_sender)
	message_from_sender_byte, err := json.Marshal(message_from_sender)
	if err != nil {
		fmt.Println("error while marshal ack - ", err)
		return
	}
	fmt.Printf(" ack - %+v\n ", message_from_sender)

	ws_conn, is_alive := Get_ws_Conn(message_from_sender.SenderID)
	if !is_alive {
		fmt.Println("websocket connection is not found")
		return
	}

	err = ws_conn.WriteMessage(msg_type, message_from_sender_byte)
	if err != nil {
		fmt.Println("error while sending ack to sender group chat- ", err)
	}
	sendMsgTo_receiver_group_chat(message_from_sender, msg_type)

}

func sendMsgTo_receiver_group_chat(message_to_receiver *models.Message, msg_type int) {

	// var Group_members []string

	Group_members, err := postgres.Get_all_group_members(message_to_receiver.GroupId, message_to_receiver.SenderDept)
	if err != nil {
		fmt.Println("error while get group members - ", err)
		return
	}
	message_to_receiver.IsAck = "not-ack"
	message_to_receiver_byte, err := json.Marshal(message_to_receiver)
	if err != nil {
		fmt.Println("error while marshal sending_msg_byte - ", message_to_receiver_byte)
		return
	}
	fmt.Println("message_to_receiver - ", message_to_receiver)

	if len(Group_members) == 0 {
		fmt.Println("group is empty !! ")
		return
	}

	for _, member := range Group_members {
		for member_id, _ := range member {
			if member_id == message_to_receiver.SenderID { //to avoid the same message send to that sender
				continue
			}

			receiver, ok := Get_ws_Conn(member_id)
			if !ok {
				fmt.Printf("%v goes to offline ! ", member_id)
				continue
			}
			err = receiver.WriteMessage(msg_type, message_to_receiver_byte)
			if err != nil {
				fmt.Println("error while sending message_to_receiver - ", err)
			}
		}

	}
}
