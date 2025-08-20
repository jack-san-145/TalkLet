package services

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
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

func FindDeptChatlistByRollno(roll_no string) string {
	dept := fmt.Sprintf("%c"+"%c", roll_no[3], roll_no[4]) // find the specific table for that email
	dept_table := dept + "_chatlist"
	return dept_table
}

func FindDeptByDept(dept string) string {
	var original_dept string
	switch dept {
	case "CSE":
		original_dept = "cs"
	case "AIDS":
		original_dept = "ad"
	case "BT":
		original_dept = "bt"
	case "ECE":
		original_dept = "ec"
	case "MECH":
		original_dept = "me"
	case "MTE":
		original_dept = "mt"
	case "IT":
		original_dept = "it"
	case "CIVIL":
		original_dept = "ce"
	case "EEE":
		original_dept = "ee"
	}
	return original_dept

}

func Find_Year_And_Section(class string) (int, string) {
	var (
		current_year int
		section      string
	)
	splitted_class := strings.Split(class, "-")
	section = splitted_class[1]
	switch splitted_class[0] {
	case "IV":
		current_year = 4
	case "III":
		current_year = 3
	case "II":
		current_year = 2
	}
	return current_year, section
}

func FindBatch(batch string) (int, int) {
	var (
		batch_year   int
		passing_year int
	)
	splitted_batch := strings.Split(batch, "-")
	batch_year, _ = strconv.Atoi(splitted_batch[0])
	passing_year, _ = strconv.Atoi(splitted_batch[1])
	return batch_year, passing_year
}

func Find_staff_or_student_by_email(email string) string {
	// Student regex
	studentRe := regexp.MustCompile(`^(\d{2})([a-z]+)(\d+)@kamarajengg\.edu\.in$`)
	// Staff regex
	staffRe := regexp.MustCompile(`^([a-z]+)([a-z]{2,3})@kamarajengg\.edu\.in$`)
	var whose_email string
	if studentRe.MatchString(email) {
		whose_email = "Student"
	}

	if staffRe.MatchString(email) {
		whose_email = "Staff"
	}
	return whose_email
}
