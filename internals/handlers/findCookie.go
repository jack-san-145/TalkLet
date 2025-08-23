package handlers

import (
	"fmt"
	"net/http"
	"tet/internals/storage/postgres"
)

func FindCookie(r *http.Request) (bool, string) {
	cookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		fmt.Println("no cookie found - ", err)
		return false, ""
	}
	roll_no, _ := postgres.FindSessionPdb(cookie.Value)
	return true, roll_no
}
