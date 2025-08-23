package postgres

import (
	"context"
	"fmt"
)

func AddNewDepartment(dept_name string) {
	addStudentTable(dept_name)
	createChatlist(dept_name)
}

func addStudentTable(dept_name string) {
	table_name := dept_name + "_students"
	unique_roll_email := dept_name + "_students_roll_email_unique"

	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
        roll_no       VARCHAR(10) PRIMARY KEY,
        register_no   TEXT UNIQUE,
        name          TEXT,
        dob           DATE,
        email         TEXT ,
        password      TEXT,
        batch_year    INT,
        passing_year  INT,
        branch        TEXT,
        current_year  INT,
        section       VARCHAR(1),
        chairperson   TEXT,
        mentor        TEXT,
		CONSTRAINT %s UNIQUE (roll_no, email)
   	  );`, table_name, unique_roll_email)
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
        group_id    SERIAL PRIMARY KEY,
        name        TEXT,
        admin       JSONB NOT NULL DEFAULT '[]'::jsonb,
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
        group_id    INT NOT NULL,
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
	arr := []string{"cs", "ec", "bt", "ad", "mt", "me", "it"} //"civil" , "eee" not added
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

func createChatlist(dept_name string) {
	table_name := dept_name + "_chatlist"
	dept := dept_name + "_students"
	query := fmt.Sprintf(`create table if not exists %s(
	sender_id varchar(10),
	receiver_id varchar(10),
	is_group boolean default false,
	group_id int,
	last_msg text,
	last_msg_id bigint,
	first_msg_id bigint,
	created_at timestamp default current_timestamp,
	primary key(sender_id,receiver_id),
	foreign key (sender_id) references %s(roll_no) on delete cascade);`, table_name, dept)
	_, err := pool.Exec(context.Background(), query)
	if err != nil {
		fmt.Println("error while creating the chatlist - ", err)
	}
}

func DropChatlistTable(dept string) {
	table_name := dept + "_chatlist"
	query := fmt.Sprintf(`drop table %s`, table_name)
	_, err := pool.Exec(context.Background(), query)
	if err != nil {
		fmt.Println("error while drop chatlist table - ", err)
	}
}
