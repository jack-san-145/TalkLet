package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"tet/internals/models"
	"tet/internals/storage/postgres"

	"golang.org/x/crypto/bcrypt"
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
	pass := r.FormValue("password")
	user.Email = r.FormValue("email")
	received_otp := r.FormValue("otp")
	user.Password, err = bcrypt.GenerateFromPassword([]byte(pass), 14)
	if err != nil {
		fmt.Println("error while generating the password - ", err)
		return
	}
	fmt.Println("Bcrypt_Password - ", user.Password)

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

func StaffRegistration(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("error while parsing the register form - ", err)
		return
	}
	var new_staff models.StaffDetails
	new_staff.StaffID = r.FormValue("staff_id")
	new_staff.Name = r.FormValue("name")
	new_staff.DOB = r.FormValue("dob")
	new_staff.Email = r.FormValue("email")
	password := r.FormValue("password")
	new_staff.Dept = r.FormValue("department")
	otp := r.FormValue("otp")
	if otp == "" {
		w.Write([]byte("<p> invalid OTP ❌<p>"))
		return
	}
	isValid := VerifyOTP_for_staff(w, r, new_staff.Email, otp)
	if !isValid {
		w.Write([]byte("<p> invalid OTP ❌<p>"))
		return
	}

	if postgres.ValidateEmail(new_staff.Email, "all_staffs") {
		w.Write([]byte("<p>You have already registered ❌</p>"))
		return
	}

	//bcrypt the staff's password
	bycrpted_password, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		fmt.Println("error while bcryp the password - ", err)
		return
	}

	new_staff.Password = string(bycrpted_password)

	fmt.Printf("new staff - %+v", new_staff)
	postgres.NewStaffRegisterPDB(new_staff)
	mutex.Lock()
	delete(OTPs, new_staff.Email)
	mutex.Unlock()
}
