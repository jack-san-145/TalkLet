package postgres

import (
	"fmt"
	"tet/internals/models"
)

func InsertToUsers(user models.User) {
	query := "insert into users(user_name,mobile_no,location,password,email) values($1,$2,$3,$4,$5)"
	_, err := Db.Exec(query, user.Name, user.Mobile_no, user.Location, user.Password, user.Email)
	if err != nil {
		fmt.Println("error while inserting user records", err)
		return
	}
	fmt.Println("users records successfuly inserted")
}
