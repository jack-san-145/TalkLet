package postgres

import (
	"context"
	"fmt"
)

func AddNewDepartment(dept_name string) {
	// addStudentTable(dept_name)
	// createChatlist(dept_name)
	CreateGroupMessageTable(dept_name)
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
    group_id         TEXT PRIMARY KEY,
    group_serial_no  SERIAL,
    name             TEXT NOT NULL DEFAULT '',
    admin            JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at       TIMESTAMP DEFAULT current_timestamp
);
`, table_name)
	_, err := pool.Exec(context.Background(), query)
	if err != nil {
		fmt.Println("error while creating the department_all_group_table - ", err)
	}
	createGroupMembersTable(dept_name, table_name)
}

func createGroupMembersTable(dept_name string, dept__group_table string) {
	table_name := dept_name + "_group_members"

	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
        group_id    TEXT NOT NULL default '',
        member_id   TEXT NOT NULL default '',
        group_name  TEXT NOT NULL default '',
        isadmin     BOOLEAN NOT NULL default false,
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
			fmt.Println("error while deleting _students", err)
		}

		query = fmt.Sprintf(`Drop table %s`, table_name+"_group_members")
		_, err = pool.Exec(context.Background(), query)
		if err != nil {
			fmt.Println("error while deleting _group_members", err)
		}
		query = fmt.Sprintf(`Drop table %s`, table_name+"_all_groups")
		_, err = pool.Exec(context.Background(), query)
		if err != nil {
			fmt.Println("error while deleting _all_groups", err)
		}
		query = fmt.Sprintf(`Drop table %s`, table_name+"_chatlist")
		_, err = pool.Exec(context.Background(), query)
		if err != nil {
			fmt.Println("error while deleting _chatlist", err)
		}
		query = fmt.Sprintf(`Drop table %s`, table_name+"group_all_messages")
		_, err = pool.Exec(context.Background(), query)
		if err != nil {
			fmt.Println("error while deleting group_all_messages", err)
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
	receiver_id varchar(10) not null default '',
	is_group boolean not null default false,
	group_id text not null default '',
	last_msg text not null default '',
	last_msg_id bigint not null default 0,
	first_msg_id bigint not null default 0,
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

func CreateGroupMessageTable(dept string) {
	table_name := dept + "_group_all_messages"
	query_for_group_message_table := fmt.Sprintf(`CREATE TABLE %s (
				msg_id BIGSERIAL ,
				sender_id VARCHAR(10) NOT NULL DEFAULT '',
				receiver_id VARCHAR(10) NOT NULL DEFAULT '',
				group_id VARCHAR(10) NOT NULL DEFAULT '',
				type TEXT NOT NULL DEFAULT '',
				content TEXT NOT NULL DEFAULT '',
				meta_data JSONB NOT NULL DEFAULT '{}'::jsonb,
				created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
				status TEXT NOT NULL DEFAULT 'sent'
			) PARTITION BY RANGE (created_at);`, table_name)
	_, err := pool.Exec(context.Background(), query_for_group_message_table)
	if err != nil {
		fmt.Println("error while creating the group message table - ", err)
		return
	}
}
