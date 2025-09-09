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

// func ValidateUser(username string, emailOrPassword string) bool {
// 	var isValid bool

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	query := "select exists(select 1 from users where user_name = $1 or email = $2)"
// 	pool.QueryRow(ctx, query, username, emailOrPassword).Scan(&isValid)
// 	fmt.Println("isvalid - ", isValid)
// 	return isValid
// }

func ValidateEmail(email string, dept_table string) bool {
	//check if the email is present or not in the corresponding table
	var isValid bool

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := fmt.Sprintf(`select exists(select 1 from %s where email = $1)`, dept_table)
	pool.QueryRow(ctx, query, email).Scan(&isValid)
	fmt.Println("isvalid - ", isValid)
	return isValid
}

func ValidateStudentLogin(roll_no string, password string) (string, bool) {
	var (
		isValid bool
	)
	_, dept_table, _ := services.Find_dept_from_rollNo(roll_no)
	column_name := "roll_no"
	isValid = isPasswordMatching(roll_no, password, dept_table, column_name)
	return roll_no, isValid
}

func ValidateStaffLogin(staff_id string, password string) (string, bool) {
	var (
		isValid bool
	)
	dept_table := "all_staffs"
	column_name := "staff_id"
	isValid = isPasswordMatching(staff_id, password, dept_table, column_name)
	return staff_id, isValid
}

func FindContact(userId string) (string, string, string, error) {
	var (
		name       string
		password   string
		email      string
		TABLE_NAME string
		SEARCH_BY  string
	)

	is_staff_or_student := services.Find_staff_or_student_by_id(userId)
	if is_staff_or_student == "STAFF" {
		TABLE_NAME = "all_staffs"
		SEARCH_BY = "staff_id"
	} else if is_staff_or_student == "STUDENT" {
		_, TABLE_NAME, _ = services.Find_dept_from_rollNo(userId)
		SEARCH_BY = "roll_no"
	}

	query := fmt.Sprintf(`select name,email,password from %s where %s = $1`, TABLE_NAME, SEARCH_BY)
	err := pool.QueryRow(context.Background(), query, userId).Scan(&name, &email, &password)
	if err == sql.ErrNoRows {
		fmt.Println("invalid user id - ", err)
		return "", "", "", fmt.Errorf("no user found")
	}
	return name, email, password, nil
}

func isPasswordMatching(roll_no string, password string, dept_table string, column_name string) bool {
	var Db_password string
	fmt.Println("your password - ", password)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	fmt.Println("dept_table - ", dept_table, " column_name - ", column_name, "roll_no - ", roll_no)
	query := fmt.Sprintf(`select password from %s where %s = $1`, dept_table, column_name)
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

func Verify_Staff(staff_id string) bool {
	var isStaff bool
	query := "select exists(select 1 from all_staffs where staff_id = $1)"
	err := pool.QueryRow(context.Background(), query, staff_id).Scan(&isStaff)
	if err != nil {
		fmt.Println("error while find the existance of the staff_id - ", err)
	}
	return isStaff
}

// used to find the group_name for both the student and staff's chatlist
func Find_groupname_by_groupid(group_id string) string {
	var group_name string
	dept := fmt.Sprintf(`%c%c`, group_id[0], group_id[1]) // here finding the dept from group_id (cs_1) -> cs
	TABLE_NAME := dept + "_all_groups"
	query := fmt.Sprintf(`select name from %s where group_id = $1`, TABLE_NAME)
	err := pool.QueryRow(context.Background(), query, group_id).Scan(&group_name)
	if err != nil {
		fmt.Println("error while accesing the group_name from group_id - ", err)
	}
	return group_name

}
