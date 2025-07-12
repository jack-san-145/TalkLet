package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
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
	query := "select user_id from users where user_name = $1 and password = $2"
	err := Db.QueryRow(query, username, password).Scan(&userId)
	if err == sql.ErrNoRows {
		isValid = false
	} else {
		isValid = true
	}
	fmt.Println("Login status - ", isValid)
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
