package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"tet/internals/services"
	"time"

	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
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

	dept_table := services.FindDeptStudentByEmail(email)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := fmt.Sprintf(`select exists(select 1 from %s where email = $1)`, dept_table)
	pool.QueryRow(ctx, query, email).Scan(&isValid)
	fmt.Println("isvalid - ", isValid)
	return isValid
}

func ValidateLogin(roll_no string, password string) (string, bool) {
	var (
		isValid bool
	)
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// query := "select user_id from users where user_name = $1"
	// err := pool.QueryRow(ctx, query, roll_no).Scan(&userId)
	// if err == sql.ErrNoRows {
	// 	fmt.Println("user not found - ", err)
	// } else {
	// 	isValid = isPasswordMatching(userId, password)
	dept_table := services.FindDeptStudentByRollNo(roll_no)
	// }
	isValid = isPasswordMatching(roll_no, password, dept_table)
	return roll_no, isValid

}

func FindUser(userId string) (string, string, string, error) {
	var (
		name     string
		password string
		email    string
	)
	dept_table := services.FindDeptStudentByRollNo(userId)

	query := fmt.Sprintf(`select name,email,password from %s where roll_no = $1`, dept_table)
	err := pool.QueryRow(context.Background(), query, userId).Scan(&name, &email, &password)
	if err == sql.ErrNoRows {
		fmt.Println("invalid user id - ", err)
		return "", "", "", fmt.Errorf("no user found")
	}
	return name, email, password, nil
}

func isPasswordMatching(roll_no string, password string, dept_table string) bool {
	var Db_password string
	fmt.Println("your password - ", password)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := fmt.Sprintf(`select password from %s where roll_no = $1`, dept_table)
	err := pool.QueryRow(ctx, query, roll_no).Scan(&Db_password)
	if err == pgx.ErrNoRows {
		fmt.Println("no roll_no found - ")
		return false
	} else if err != nil {
		fmt.Println("error while accessing the db password - ", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(Db_password), []byte(password))
	if err != nil {
		return false
	} else {
		return true
	}
}
