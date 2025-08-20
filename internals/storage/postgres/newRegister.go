package postgres

import (
	"context"
	"fmt"
	"tet/internals/models"
	"time"
)

func InsertToUsers(user models.User) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "insert into users(user_name,mobile_no,location,password,email) values($1,$2,$3,$4,$5)"
	_, err := pool.Exec(ctx, query, user.Name, user.Mobile_no, user.Location, string(user.Password), user.Email)
	if err != nil {
		fmt.Println("error while inserting user records", err)
		return
	}
	fmt.Println("users records successfuly inserted")
}

func NewStaffRegisterPDB(new_staff models.StaffDetails) {
	query := "insert into all_staffs(staff_id,name,dob,email,password,dept) values($1,$2,$3,$4,$5,$6)"
	_, err := pool.Exec(context.Background(), query,
		new_staff.StaffID,
		new_staff.Name,
		new_staff.DOB,
		new_staff.Email,
		new_staff.Password,
		new_staff.Dept)
	if err != nil {
		fmt.Println("error while inserting the staff registration details - ", err)
	}
}
