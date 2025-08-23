package postgres

import (
	"context"
	"fmt"
	"tet/internals/models"
	"tet/internals/services"
	"time"
)

func CreateNewGroupPDB(admin string, group_name string, department string, Students []models.StudentDetails) {

	insert_student_into_dept_table(department, Students)
	group_id, _, err := add_new_group(admin, group_name, department)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	add_group_members(department, group_id, group_name, admin, Students)

}

func insert_student_into_dept_table(department string, Students []models.StudentDetails) {
	fmt.Println("students - ", Students)
	dept_table := services.FindDeptByDept(department) + "_students"
	query := fmt.Sprintf(`insert into %s(roll_no,register_no,name,dob,
		email,batch_year,passing_year,branch,current_year,
		section,chairperson,mentor) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		on conflict (roll_no,email) do nothing`, dept_table)

	for _, new_student := range Students {
		fmt.Printf("name - %v date from pdb - %v\n ", new_student.Name, new_student.Email)
		_, err := pool.Exec(context.Background(), query,
			new_student.Roll_no,
			new_student.Register_no,
			new_student.Name,
			new_student.DOB,
			new_student.Email,
			new_student.Batch_year,
			new_student.Passing_year,
			new_student.Branch,
			new_student.Current_year,
			new_student.Section,
			new_student.Chairperson,
			new_student.Mentor)
		if err != nil {
			fmt.Println("error while inserting new students to their dept table - ", err)
		}
	}
}

func add_new_group(admin string, group_name string, department string) (int32, string, error) {
	var group_id int32

	//adding new group to the dept_all_groups and returning its group_id
	dept_table := services.FindDeptByDept(department) + "_all_groups"
	created_at := time.Now().Format("2006-01-02 15:04:05")
	// admin = "T2505778"
	adminJson := fmt.Sprintf(`["%s"]`, admin)
	query := fmt.Sprintf(`insert into %s(name,admin,created_at) values($1,$2::jsonb ,$3) returning group_id`, dept_table)
	err := pool.QueryRow(context.Background(), query, group_name, adminJson, created_at).Scan(&group_id)
	if err != nil {
		actual_err := fmt.Sprint("error while adding new group - ", err)
		return 0, "", fmt.Errorf(actual_err)
	}
	return group_id, dept_table, nil
}

func add_group_members(dept string, group_id int32, group_name string, admin string, students []models.StudentDetails) {
	table_name := services.FindDeptByDept(dept) + "_group_members"
	//query to add only the admin by their staff_id
	query := fmt.Sprintf(`insert into %s(group_id,member_id,group_name,isadmin) values($1,$2,$3,$4)`, table_name)
	_, err := pool.Exec(context.Background(), query, group_id, admin, group_name, true)
	if err != nil {
		fmt.Println("error while add admin to new group - ", err)
	}

	//query to all the students of this group by their roll_no

	for _, new_student := range students {
		query = fmt.Sprintf(`insert into %s(group_id,member_id,group_name,isadmin) values($1,$2,$3,$4)`, table_name)
		_, err := pool.Exec(context.Background(), query, group_id, new_student.Roll_no, group_name, false)
		if err != nil {
			fmt.Println("error while add students to new group - ", err)
		}
	}

}
