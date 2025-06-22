package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"tet/internals/models"
	"tet/internals/storage/postgres"
)

func ServeRegister(w http.ResponseWriter, r *http.Request) {
	templ, err := template.ParseFiles("../../ui/templates/register.html")
	if err != nil {
		fmt.Println("register not found")
		return
	}
	err = templ.Execute(w, nil)
}

func AccountRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := r.ParseForm()
	if err != nil {
		fmt.Println("error while parsing account register form")
		return
	}

	user.Name = r.FormValue("username")
	user.Mobile_no = r.FormValue("mobile_no")
	user.Location = r.FormValue("location")
	user.Password = r.FormValue("password")
	user.Email = r.FormValue("email")
	received_otp := r.FormValue("otp")

	mutex.Lock()
	sent_otp := OTPs[user.Email]
	delete(OTPs, user.Email)
	mutex.Unlock()

	fmt.Println("received_otp - ", received_otp)
	fmt.Println("sent_otp - ", sent_otp)
	if sent_otp != received_otp {
		w.Write([]byte("<p>OTP invalid ❌<p>"))
	}
	fmt.Println(user)
	postgres.InsertToUsers(user)
	w.Header().Set("Hx-Redirect", "/")
}
