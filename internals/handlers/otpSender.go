package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"regexp"
	"strconv"
	"strings"
	"tet/internals/services"
	"tet/internals/storage/postgres"
	"tet/internals/storage/redis"
	"time"
)

func SendOtpHandler_for_students(w http.ResponseWriter, r *http.Request) {

	fmt.Println("body - ", r.Body)
	err := r.ParseForm()
	if err != nil {
		fmt.Println("error while parsing form")
		return
	}
	email := r.FormValue("email")
	fmt.Println("email - ", email)
	isMatch, _ := regexp.MatchString(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`, email)
	if !isMatch {
		w.Write([]byte("<p>Invalid Email ❌</p>"))
		return
	}
	isGmail := strings.Split(email, "@")
	if isGmail[1] != "kamarajengg.edu.in" {
		w.Write([]byte("<p>Invalid Email ❌</p>"))
		return
	}

	if !postgres.ValidateEmail(email, services.FindDeptStudentByEmail(email)) {
		w.Write([]byte("<p>You have not registered yet ❌</p>"))
		return
	}
	otp := generateOtp()
	go sendEmailTo(email, otp)
	w.Write([]byte("<p>✅ Otp sent to " + email))

	redis.Set_OTP_to_redis(email, otp) //store otp to  redis

}

func SendOtpHandler_for_staffs(w http.ResponseWriter, r *http.Request) {

	fmt.Println("body - ", r.Body)
	err := r.ParseForm()
	if err != nil {
		fmt.Println("error while parsing form")
		return
	}
	email := r.FormValue("email")
	fmt.Println("email - ", email)

	// to verify a email for either staff or not
	email = strings.TrimSpace(email)
	isStaffMail := services.Find_staff_or_student_by_email(email)
	if isStaffMail != "Staff" {
		fmt.Println("isStaffMail - ", isStaffMail)
		w.Write([]byte("<p>You are not a staff ❌</p>"))
		return
	}

	if postgres.ValidateEmail(email, "all_staffs") {
		w.Write([]byte("<p>You have already registered ❌</p>"))
		// WriteJSON(w, r, map[string]string{"status": "You have already registered ❌"})
		return
	}
	otp := generateOtp()
	go sendEmailTo(email, otp)
	w.Write([]byte("<p>✅ Otp sent to " + email))

	redis.Set_OTP_to_redis(email, otp) //set otp to redis

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

func VerifyOTP_for_student_handler(w http.ResponseWriter, r *http.Request) {
	var otp_response = make(map[string]bool)
	err := r.ParseForm()
	if err != nil {
		fmt.Println("error while parsing otp form - ", err)
		return
	}
	received_email := r.FormValue("email")
	received_otp := r.FormValue("otp")

	sent_otp, err := redis.Get_OTP_from_redis(received_email)
	if err != nil {
		otp_response["otp_valid"] = false
	}

	if sent_otp != received_otp {
		// w.Write([]byte("<p> invalid OTP ❌<p>"))
		otp_response["otp_valid"] = false
	} else {
		otp_response["otp_valid"] = true
	}
	WriteJSON(w, r, otp_response)
	// w.Header().Set("Hx-Redirect", "/")

}
func VerifyOTP_for_staff_handler(w http.ResponseWriter, r *http.Request) {
	var otp_response = make(map[string]bool)
	err := r.ParseForm()
	if err != nil {
		fmt.Println("error while parsing staff otp form - ", err)
		return
	}
	received_email := r.FormValue("email")
	received_otp := r.FormValue("otp")

	sent_otp, err := redis.Get_OTP_from_redis(received_email)
	if err != nil {
		otp_response["otp_valid"] = false
	}

	if sent_otp != received_otp {
		// w.Write([]byte("<p> invalid OTP ❌<p>"))
		otp_response["otp_valid"] = false
	} else {
		otp_response["otp_valid"] = true

	}
	WriteJSON(w, r, otp_response)

}

func VerifyOTP_for_staff(w http.ResponseWriter, r *http.Request, received_email string, received_otp string) bool {

	sent_otp, err := redis.Get_OTP_from_redis(received_email)
	if err != nil {
		return false
	}

	if sent_otp != received_otp {
		return false
	}
	return true

}
