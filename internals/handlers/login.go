package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"tet/internals/storage/postgres"
	"tet/internals/storage/redis"

	"time"
)

func ServeLogin(w http.ResponseWriter, r *http.Request) {
	templ, err := template.ParseFiles("../../ui/templates/login.html")
	if err != nil {
		fmt.Println("login not found")
		return
	}
	err = templ.Execute(w, nil)
	if err != nil {
		fmt.Println("login not serve")
	}
}

func StudentLoginValidationHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("error while parsing login")
		return
	}
	roll_no := r.FormValue("roll_no")
	password := r.FormValue("password")
	roll_no, isValid := postgres.ValidateStudentLogin(roll_no, password)
	if !isValid {
		w.Write([]byte("<p>Invalid username or Password ❌</p>"))
		return
	}
	createSession(roll_no, w)
}

func StaffLoginValidationHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("error while parsing login")
		return
	}
	staff_id := r.FormValue("staff_id")
	password := r.FormValue("password")
	staff_id, isValid := postgres.ValidateStaffLogin(staff_id, password)
	if !isValid {
		w.Write([]byte("<p>Invalid username or Password ❌</p>"))
		return
	}
	createSession(staff_id, w)
}

func createSession(id string, w http.ResponseWriter) {
	session := postgres.GenerateSessionID(id)
	redis.SetSessionToRdb(session)
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.Session_id,
		Path:     "/",
		Expires:  time.Now().Add(3 * time.Hour),
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		// SameSite: http.SameSiteStrictMode,
	})
	fmt.Println("successfully logged in")
	w.Header().Set("Hx-Redirect", "/talklet/serve-index")
}
