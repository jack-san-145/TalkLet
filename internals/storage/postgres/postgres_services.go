package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"tet/internals/services"
)

func Find_dept_from_staff_id(staff_id string) (string, error) {
	var (
		dept_from_db  string
		original_dept string
	)

	//get the dept from staff table -> CSE
	query := "select dept from all_staffs where staff_id = $1"
	err := pool.QueryRow(context.Background(), query, staff_id).Scan(&dept_from_db)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("invalid staff id")
	} else if err == nil {
		fmt.Println("error while finding dept from all_staffs table - ", err)
		return "", fmt.Errorf(err.Error())
	}

	//get the dept from dept -> (CSE -> cs )
	original_dept = services.FindDeptByDept(dept_from_db)
	return original_dept, nil
}
