package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tet/internals/models"
	"tet/internals/storage/postgres"

	"golang.org/x/crypto/bcrypt"
)

func SetPassword(w http.ResponseWriter, r *http.Request) {

	var (
		bcryted_password []byte
		user_password    models.Password
	)
	err := json.NewDecoder(r.Body).Decode(&user_password)
	if err != nil {
		fmt.Println("error while decoding the password - ", err)
		return
	}
	fmt.Println("user_pass - ", user_password.Pass)

	bcryted_password, err = bcrypt.GenerateFromPassword([]byte(user_password.Pass), 14)
	if err != nil {
		fmt.Println("error while generating the password - ", err)
		return
	}
	fmt.Println("Bcrypt_Password - ", string(bcryted_password))
	err = postgres.SetPasswordDB(user_password.Email, string(bcryted_password))
	fmt.Println("error - ", err)
	if err == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

}
