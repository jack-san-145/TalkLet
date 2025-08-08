package services

import (
	"fmt"
)

func FindDeptStudentByEmail(email string) string {
	//23ucs145@kamarajengg.edu.in
	dept := fmt.Sprintf("%c"+"%c", email[3], email[4]) // find the specific table for that email
	dept_table := dept + "_students"
	return dept_table
}

func FindDeptStudentByRollNo(roll_no string) string {
	dept := fmt.Sprintf("%c"+"%c", roll_no[3], roll_no[4]) // find the specific table for that email
	dept_table := dept + "_students"
	return dept_table
}
