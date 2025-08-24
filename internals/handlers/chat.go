package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"tet/internals/storage/postgres"
)

func Chatlist(w http.ResponseWriter, r *http.Request) {
	isFound, userId := FindCookie(r)
	if !isFound {
		return
	}
	Chatlist := postgres.LoadChatlist(userId)
	fmt.Println(Chatlist)
	WriteJSON(w, r, Chatlist)
}

func LoadChatMessages(w http.ResponseWriter, r *http.Request) {
	isFound, userId := FindCookie(r)
	if !isFound {
		return
	}
	// query paramater = /talklet/chat-history/${contact_id}?limit=${limit}&offset=${offset}
	contactID := chi.URLParam(r, "contact_id")
	limit := r.URL.Query().Get("limit") // to get the limit value from the request
	offset := r.URL.Query().Get("offset")

	limit_int, _ := strconv.Atoi(limit)
	offset_int, _ := strconv.Atoi(offset)
	fmt.Println("contactID - ", contactID) //means the receiver's roll no
	AllMessages, err := postgres.LoadOTOChatMessagesPDb(userId, contactID, limit_int, offset_int)
	if err != nil {
		WriteJSON(w, r, "the chat is empty")
		return
	}
	WriteJSON(w, r, AllMessages)

}
