package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tet/internals/handlers"
	"tet/internals/models"
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

	var (
		Send_msg_To       models.Message
		received_msg_From models.Message
	)

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	websocketConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("error while upgrade http to websocket - ", err)
	}
	Set_ws_conn(senderID, websocketConn) //store the new ws connection to the ConnMap

	ConnMap.Range(func(key any, value any) bool {
		id, ok := key.(string)
		// conn:=value.(*websocket.Conn)
		fmt.Printf("id - %v is connected\n", id)
		return ok
	})
	Get_ws_Conn(senderID)
	defer func() {
		del_ws_conn(senderID) //to delete the ws_connection from the ConnMap
	}()

	for {

		//receiving message
		msg_type, msg, err := websocketConn.ReadMessage()
		msg_time := time.Now().Format("2006-01-02 15:04:05")
		received_msg_From.CreatedAt = msg_time
		received_msg_From.SenderID = senderID
		if err != nil {
			fmt.Println("error while reading ws- ", err)
			return
		}
		err = json.Unmarshal(msg, &received_msg_From)
		if err != nil {
			fmt.Println("error while unmarshal - ", err)
		}

		sendAck(&received_msg_From, msg_type)
		fmt.Println("received message  details after ack - ", received_msg_From)

		// sending the received msg
		Send_msg_To.ID = received_msg_From.ID
		Send_msg_To.SenderID = senderID
		Send_msg_To.ReceiverID = received_msg_From.ReceiverID
		Send_msg_To.Content = received_msg_From.Content
		Send_msg_To.Type = received_msg_From.Type
		Send_msg_To.CreatedAt = msg_time
		sendMsgTo(&Send_msg_To, msg_type) //send to frd

		// for _, conn := range ConnMap {
		// 	if conn == websocketConn {
		// 		continue
		// 	}
		// 	// temp := []byte(fmt.Sprintf("from client - %v", size))
		// 	// combined_msg := append(temp, msg...)
		// 	sending_data, _ := json.Marshal(Send_msg_To)
		// 	err := conn.WriteMessage(msg_type, sending_data)
		// 	if err != nil {
		// 		fmt.Println("error while write ws - ", err)
		// 		delete(ConnMap, senderID)
		// 		return
		// 	}
		// 	fmt.Println("sending message to details - ", Send_msg_To)

		// }
	}
}

func sendAck(msg_sent_by_sender *models.Message, msg_type int) {
	msg_sent_by_sender.IsAck = "ack"

	// var temp models.Message
	if msg_type == 1 { // which is websocket.TextMessage(1)
		msg_sent_by_sender.Type = "text/plain"
		msg_sent_by_sender.ID = postgres.Store_Privatechat_MessagesPostDB(*msg_sent_by_sender)
		msg_sent_by_sender.Status = "sent"
	}

	go postgres.AddLastMsgToChatlist(msg_sent_by_sender.SenderID, msg_sent_by_sender.ReceiverID, msg_sent_by_sender.ID, msg_sent_by_sender.Content, msg_sent_by_sender.CreatedAt)
	msg_sent_by_sender_byte, err := json.Marshal(msg_sent_by_sender)
	if err != nil {
		fmt.Println("error while marshal ack - ", err)
		return
	}
	fmt.Printf(" ack - %+v\n ", msg_sent_by_sender)

	ws_conn, is_alive := Get_ws_Conn(msg_sent_by_sender.SenderID)
	if !is_alive {
		fmt.Println("websocket connection is not found")
		return
	}

	err = ws_conn.WriteMessage(msg_type, msg_sent_by_sender_byte)
	if err != nil {
		fmt.Println("error while sending ack to sender - ", err)
	}

}

func sendMsgTo(msg_send_to_frd *models.Message, msg_type int) {
	msg_send_to_frd.IsAck = "not-ack"
	sending_msg_byte, err := json.Marshal(msg_send_to_frd)
	if err != nil {
		fmt.Println("error while marshal sending_msg_byte - ", sending_msg_byte)
		return
	}
	fmt.Println("msg_send_to_frd - ", msg_send_to_frd)

	receiver_id, ok := Get_ws_Conn(msg_send_to_frd.ReceiverID)
	if !ok {
		fmt.Printf("%v goes to offline - ", msg_send_to_frd.ReceiverID)
		return
	}
	err = receiver_id.WriteMessage(msg_type, sending_msg_byte)
	if err != nil {
		fmt.Println("error while sending msg_send_to_frd - ", err)
	}
}
