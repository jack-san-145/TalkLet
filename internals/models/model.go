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
	Roll_no    string
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
	ID         int64     `json:"msg_id"`
	SenderID   string    `json:"sender_id"`
	ReceiverID string    `json:"receiver_id"`
	Type       string    `json:"type"`
	Content    string    `json:"content"`
	FileURL    string    `json:"file_url"`
	FileSize   int64     `json:"file_size"`
	MimeType   string    `json:"mime_type"`
	CreatedAt  string `json:"created_at"`
	IsAck      string    `json:"is_ack"`
	Status     string    `json:"status"`
}

type ChatlistToSend struct {
	UserId      string `json:"user_id"`
	ContactId   string `json:"contact_id"`
	ContactName string `json:"contact_name"`
	LastMsg     string `json:"last_msg"`
	CreatedAt   string `json:"created_at"`
}

type ChatlistForLocal struct {
	UserID     string
	ContactId  string
	IsGroup    bool
	GroupId    int32
	LastMsg    string
	LastMsgId  int64
	FirstMsgId int64
	CreatedAt  time.Time
}

type NewGroup struct {
	Admin   []string `json:"group_admin"`
	Name    string   `json:"group_name"`
	Members []string `json:"group_members"`
}
