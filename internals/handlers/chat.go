package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"tet/internals/services"
	"tet/internals/storage/postgres"

	"github.com/go-chi/chi/v5"
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

func GetAllContactsHandeler(w http.ResponseWriter, r *http.Request) {
	isFound, userID := FindCookie(r)
	if !isFound {
		return
	}

	_, dpt_table, _ := services.Find_dept_from_rollNo(userID)
	Allcontacts := postgres.GetAllContacts(dpt_table)
	WriteJSON(w, r, Allcontacts)

}
