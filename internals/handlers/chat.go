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

func LoadPrivateChatMessages(w http.ResponseWriter, r *http.Request) {
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
	AllMessages, err := postgres.Load_PrivateChatMessages_PDb(userId, contactID, limit_int, offset_int)
	if err != nil {
		WriteJSON(w, r, "the chat is empty")
		return
	}
	fmt.Println("from the LoadChatMessages handlers - ", AllMessages)
	WriteJSON(w, r, AllMessages)

}

func LoadGroupChatMessages(w http.ResponseWriter, r *http.Request) {
	isFound, userId := FindCookie(r)
	if !isFound {
		return
	}
	// query paramater = /talklet/chat-history/${contact_id}?limit=${limit}&offset=${offset}
	groupId := chi.URLParam(r, "group_id")
	limit := r.URL.Query().Get("limit") // to get the limit value from the request
	offset := r.URL.Query().Get("offset")

	limit_int, _ := strconv.Atoi(limit)
	offset_int, _ := strconv.Atoi(offset)
	fmt.Println("groupId - ", groupId) //means the receiver's roll no
	AllMessages, err := postgres.Load_GroupChatMessages_PDb(userId, groupId, limit_int, offset_int)
	if err != nil {
		WriteJSON(w, r, "the group chat is empty")
		return
	}
	fmt.Println("from the LoadChatMessages handlers - ", AllMessages)
	WriteJSON(w, r, AllMessages)

}
