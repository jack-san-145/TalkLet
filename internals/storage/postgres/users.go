package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func ValidateUser(username string, emailOrPassword string) bool {
	var isValid bool
	query := "select exists(select 1 from users where user_name = $1 or email = $2)"
	Db.QueryRow(query, username, emailOrPassword).Scan(&isValid)
	fmt.Println("isvalid - ", isValid)
	return isValid
}

func ValidateLogin(username string, password string) (int, bool) {
	var (
		userId  int
		isValid bool
	)
	query := "select user_id from users where user_name = $1"
	err := Db.QueryRow(query, username).Scan(&userId)
	if err == sql.ErrNoRows {
		fmt.Println("user not found - ", err)
	} else {
		isValid = isPasswordMatching(userId, password)
	}
	return userId, isValid

}

func FindUser(userId int) (int, string, string, string, string, string, error) {
	var (
		userName string
		mobileNo string
		location string
		password string
		email    string
	)
	query := "select * from users where user_id = $1 "
	err := Db.QueryRow(query, userId).Scan(&userId, &userName, &mobileNo, &location, &password, &email)
	if err == sql.ErrNoRows {
		fmt.Println("invalid user id - ", err)
		return 0, "", "", "", "", "", fmt.Errorf("no user found")
	}
	return userId, userName, mobileNo, location, password, email, nil
}

func isPasswordMatching(userId int, password string) bool {
	var Db_password string
	query := "select password from users where user_id = $1"
	err := Db.QueryRow(query, userId).Scan(&Db_password)
	if err != nil {
		fmt.Println("error while accessing the db password - ", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(Db_password), []byte(password))
	if err != nil {
		return false
	} else {
		return true
	}
}
