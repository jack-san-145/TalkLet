package postgres

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func ValidateUser(username string, emailOrPassword string) bool {
	var isValid bool

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "select exists(select 1 from users where user_name = $1 or email = $2)"
	pool.QueryRow(ctx, query, username, emailOrPassword).Scan(&isValid)
	fmt.Println("isvalid - ", isValid)
	return isValid
}

func ValidateEmail(email string) bool {
	var isValid bool

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "select exists(select 1 from users where email = $1)"
	pool.QueryRow(ctx, query, email).Scan(&isValid)
	fmt.Println("isvalid - ", isValid)
	return isValid
}

func ValidateLogin(username string, password string) (int, bool) {
	var (
		userId  int
		isValid bool
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "select user_id from users where user_name = $1"
	err := pool.QueryRow(ctx, query, username).Scan(&userId)
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "select * from users where user_id = $1 "
	err := pool.QueryRow(ctx, query, userId).Scan(&userId, &userName, &mobileNo, &location, &password, &email)
	if err == sql.ErrNoRows {
		fmt.Println("invalid user id - ", err)
		return 0, "", "", "", "", "", fmt.Errorf("no user found")
	}
	return userId, userName, mobileNo, location, password, email, nil
}

func isPasswordMatching(userId int, password string) bool {
	var Db_password string

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "select password from users where user_id = $1"
	err := pool.QueryRow(ctx, query, userId).Scan(&Db_password)
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
