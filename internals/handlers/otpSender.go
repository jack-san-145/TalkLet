package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"regexp"
	"strconv"
	"sync"
	"tet/internals/storage/postgres"
	"time"
)

var OTPs = make(map[string]string)

var mutex sync.Mutex

func SendOtpRegisterHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("body - ", r.Body)
	err := r.ParseForm()
	if err != nil {
		fmt.Println("error while parsing form")
		return
	}
	username := r.FormValue("username")
	email := r.FormValue("email")
	fmt.Println("email - ", email)
	isMatch, _ := regexp.MatchString(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`, email)
	if !isMatch {
		w.Write([]byte("<p>Invalid Email ❌❌</p>"))
		return
	}
	if postgres.ValidateUser(username, email) {
		w.Write([]byte("<p>Account already exists ❌</p>"))
		return
	}
	otp := generateOtp()
	go sendEmailTo(email, otp)
	w.Write([]byte("<p>✅ Otp sent to " + email))

	mutex.Lock()
	OTPs[email] = otp
	mutex.Unlock()

}

func sendEmailTo(email string, otp string) {

	from := "talkletprivatelimited@gmail.com"
	password := os.Getenv("TALKLET_EMAIL_PASSWORD")
	fmt.Println("password send mail - ", password)
	msg := []byte("Subject: Your Otp for TalkLet - " + otp + "\r\n\r\nYour One time OTP is " + otp + "\r\n\r\n- TalkLet Team")
	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, from, []string{email}, msg)
	if err != nil {
		fmt.Println("Email not sent - ", err)
		return
	}
	fmt.Println("Email sent ")
}

func generateOtp() string {
	rand.Seed(time.Now().UnixNano())
	otp := rand.Intn(900000) + 100000
	fmt.Println("generated otp - ", otp)
	otp_string := strconv.Itoa(otp)
	return otp_string
}
