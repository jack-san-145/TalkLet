package handlers

import (
	"fmt"
	"net/http"
	"tet/internals/storage/postgres"
)

func FindCookie(r *http.Request) (bool, int) {
	cookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		fmt.Println("no cookie found ")
		return false, 0
	}
	userId, _ := postgres.FindSessionPdb(cookie.Value)
	return true, userId
}
