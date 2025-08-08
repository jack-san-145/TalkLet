package postgres

import (
	"context"
	"fmt"
)

func AddNewDepartment(dept_name string) {
	addStudentTable(dept_name)
}

func addStudentTable(dept_name string) {
	table_name := dept_name + "_students"

	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
        roll_no       VARCHAR(10) PRIMARY KEY,
        register_no   TEXT UNIQUE,
        name          TEXT,
        dob           DATE,
        email         TEXT UNIQUE,
        password      TEXT,
        batch_year    INT,
        passing_year  INT,
        branch        TEXT,
        current_year  INT,
        section       VARCHAR(1),
        chairperson   TEXT,
        mentor        TEXT
   	  );`, table_name)
	_, err := pool.Exec(context.Background(), query)
	if err != nil {
		fmt.Println("error while creating the department_student_table - ", err)
		return
	}
	createDepartmentGroupTable(dept_name)
}

func createDepartmentGroupTable(dept_name string) {
	table_name := dept_name + "_all_groups"

	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
        group_id    TEXT PRIMARY KEY,
        name        TEXT,
        admin       JSONB,
        created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    ); `, table_name)
	_, err := pool.Exec(context.Background(), query)
	if err != nil {
		fmt.Println("error while creating the department_all_group_table - ", err)
	}
	createGroupMembersTable(dept_name, table_name)
}

func createGroupMembersTable(dept_name string, dept__group_table string) {
	table_name := dept_name + "_group_members"

	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
        group_id    TEXT NOT NULL,
        member_id   TEXT NOT NULL,
        group_name  TEXT,
        isadmin     BOOLEAN,
        PRIMARY KEY (group_id, member_id),
		FOREIGN KEY (group_id) REFERENCES %s(group_id) ON DELETE CASCADE
    );`, table_name, dept__group_table)
	_, err := pool.Exec(context.Background(), query)
	if err != nil {
		fmt.Println("error while creating the department_group_memners_table - ", err)
	}
}

func DropAllTable() {
	var query string
	var err error
	arr := []string{"cse", "ece", "bt", "aids", "mtre", "eee", "mech", "civil", "it"}
	for _, table_name := range arr {
		query = fmt.Sprintf(`Drop table %s`, table_name+"_students")
		_, err = pool.Exec(context.Background(), query)
		if err != nil {
			fmt.Println("error while deleting 1st", err)
		}

		query = fmt.Sprintf(`Drop table %s`, table_name+"_group_members")
		_, err = pool.Exec(context.Background(), query)
		if err != nil {
			fmt.Println("error while deleting 3rd", err)
		}
		query = fmt.Sprintf(`Drop table %s`, table_name+"_all_groups")
		_, err = pool.Exec(context.Background(), query)
		if err != nil {
			fmt.Println("error while deleting 2nd", err)
		}
	}

}

func AlterTable() {
	var query string
	var err error
	arr := []string{"cse", "ece", "bt", "aids", "mtre", "eee", "mech", "civil", "it"}
	for _, table_name := range arr {
		query = fmt.Sprintf(`alter table %s add column public_key text default null`, table_name+"_students")
		_, err = pool.Exec(context.Background(), query)
		if err != nil {
			fmt.Println("error while altering the table rows by public key - ", err)
		}
	}

}
