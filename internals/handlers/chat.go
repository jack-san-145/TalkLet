package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"tet/internals/storage/postgres"

	"github.com/go-chi/chi/v5"
)

func OneToOneChatlist(w http.ResponseWriter, r *http.Request) {
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
	contactID := chi.URLParam(r, "contact_id")
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	contactID_int, _ := strconv.Atoi(contactID)
	limit_int, _ := strconv.Atoi(limit)
	offset_int, _ := strconv.Atoi(offset)
	fmt.Println("contactID - ", contactID_int)
	AllMessages, err := postgres.LoadChatMessagesPDb(userId, contactID_int, limit_int, offset_int)
	if err != nil {
		WriteJSON(w, r, "the chat is empty")
		return
	}
	WriteJSON(w, r, AllMessages)

}
