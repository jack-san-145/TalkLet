package postgres

import (
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

func ValidateLogin(username string, password string) bool {
	var isValid bool
	query := "select exists(select 1 from users where user_name = $1 and password = $2)"
	Db.QueryRow(query, username, password).Scan(&isValid)
	fmt.Println("Login status - ", isValid)
	return isValid
}
