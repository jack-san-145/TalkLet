package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"tet/internals/services"
)

func Find_dept_from_staff_id(staff_id string) (string, string, string, error) {

	var (
		dept                 string
		staff_table          string
		staff_chatlist_table string
	)

	query := "select dept from all_staffs where staff_id = $1"
	err := pool.QueryRow(context.Background(), query, staff_id).Scan(&dept)
	if err != nil {
		fmt.Println("error while finding the staffs department - ", err)
		return "", "", "", fmt.Errorf(err.Error())
	}
	dept = services.FindDeptByDept(dept)
	staff_table = "all_staffs"
	staff_chatlist_table = "all_staffs_chatlist"
	return dept, staff_table, staff_chatlist_table, nil
}

func Get_all_group_members(group_id string, dept string) ([]map[string]string, error) {
	var All_group_members []map[string]string
	fmt.Println("before group_id - ", group_id)
	fmt.Println("dept - ", dept)
	group_table_name := dept + "_group_members" // cs_group_members
	fmt.Println("group_table_name - ", group_table_name)
	group_members_query := fmt.Sprintf(`select member_id from %s where group_id = $1`, group_table_name)
	Rows, err := pool.Query(context.Background(), group_members_query, group_id)
	fmt.Println("after group_id - ", group_id)
	if err == sql.ErrNoRows {
		fmt.Println("group is empty ")
		return nil, fmt.Errorf("group is empty")
	} else if err != nil {
		fmt.Println("error while accessing group members - ", err)
		return nil, fmt.Errorf("err - ", err)
	}
	fmt.Println("group_members_query - ", group_members_query)
	for Rows.Next() {
		fmt.Println("comming inside the loop ")
		var member_and_dept = make(map[string]string)
		var (
			group_member_id string
			member_chatlist string
		)
		err := Rows.Scan(&group_member_id)
		if err != nil {
			fmt.Println("error while fetching the group")
			return nil, fmt.Errorf(err.Error())
		}

		staff_or_student := services.Find_staff_or_student_by_id(group_member_id)
		if staff_or_student == "STUDENT" {
			_, _, member_chatlist = services.Find_dept_from_rollNo(group_member_id)
		} else if staff_or_student == "STAFF" {
			_, _, member_chatlist, _ = Find_dept_from_staff_id(group_member_id)
		}

		member_and_dept[group_member_id] = member_chatlist

		All_group_members = append(All_group_members, member_and_dept) //adding maps with id with department to the All_group_members
	}
	fmt.Println("All_group_members - ", All_group_members)
	return All_group_members, nil
}
