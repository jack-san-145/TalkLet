package models

import (
	"time"
)

type User struct {
	Name      string
	Mobile_no string
	Location  string
	Password  string
	Email     string
}

type Session struct {
	Session_id string
	User_id    int
	Expires_at time.Time
}

type Profile struct {
	UserID int    `json:"userId"`
	Name   string `json:"name"`
	About  string `json:"about"`
	Mobile string `json:"mobile"`
	Email  string `json:"email"`
	Joined string `json:"joined"`
}

type Message struct {
	ID         int       `json:"msg_id"`
	SenderID   int       `json:"sender_id"`
	ReceiverID int       `json:"receiver_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	IsAck      string    `json:"is_ack"`
	Type       string    `json:"type"`
}

type Chatlist struct {
	ContactId   int
	ContactName string
	LastMsg     string
	CreatedAt   string
}

// type MessageDB struct {
// 	MsgId      int       `json:"msg_id"`
// 	SenderID   int       `json:"sender_id"`
// 	ReceiverID int       `json:"receiver_id"`
// 	Type       string    `json:"type"`
// 	Content    string    `json:"content"`
// 	CreatedAt  time.Time `json:"created_at"`
// }
