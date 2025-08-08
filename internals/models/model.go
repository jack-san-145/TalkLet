package models

import (
	"time"
)

type Password struct {
	Email string `json:"email"`
	Pass  string `json:"password"`
}

type User struct {
	Name      string
	Mobile_no string
	Location  string
	Password  []byte
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

type NewGroup struct {
	Admin   []int  `json:"group_admin"`
	Name    string `json:"group_name"`
	Members []int  `json:"group_members"`
}
