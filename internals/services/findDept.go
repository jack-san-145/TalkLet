package services

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func FindDeptStudentByEmail(email string) string {
	var (
		dept                string
		dept_students_table string
	)
	var emailRegex = regexp.MustCompile(`^\d{2}(u|p)(cs|ad|it|ee|ec|mt|bt|ce|me)[a-z]?(\d+)@`)
	matches := emailRegex.FindStringSubmatch(email)
	if len(matches) > 1 {
		dept = matches[2] // matches = [23ucs145@ u cs 145]
		dept_students_table = dept + "_students"
	}
	return dept_students_table
}

func Find_dept_from_rollNo(rollNo string) (string, string, string) {
	var (
		dept                string
		dept_students_table string
		dept_chatlist_table string
	)
	var rollRegex = regexp.MustCompile(`^\d{2}[up](cs|it|ee|ec|mt|bt|ce|me)(\d+)$`)

	matches := rollRegex.FindStringSubmatch(rollNo)
	if len(matches) > 1 {
		dept = matches[2] // matches = [23ucs145 u cs 145]
		dept_students_table = dept + "_students"
		dept_chatlist_table = dept + "_chatlist"
	}
	return dept, dept_students_table, dept_chatlist_table
}

func FindDeptChatlistByRollno(rollNo string) string {
	var (
		dept                string
		dept_chatlist_table string
	)
	var rollRegex = regexp.MustCompile(`^\d{2}(u|p)(cs|it|ee|ec|mt|bt|ce|me)(\d+)$`)
	matches := rollRegex.FindStringSubmatch(rollNo)
	if len(matches) > 1 {
		dept = matches[2] // matches = [23ucs145 u cs 145]

		dept_chatlist_table = dept + "_chatlist"
	}
	fmt.Println("matches in chatlist - ", matches)
	return dept_chatlist_table
}

// func FindDeptStudentByEmail(email string) string {
// 	//23ucs145@kamarajengg.edu.in
// 	dept := fmt.Sprintf("%c"+"%c", email[3], email[4]) // find the specific table for that email
// 	dept_table := dept + "_students"
// 	return dept_table
// }

// func Find_dept_from_rollNo(roll_no string) (string, string, string) {
// 	//23ucs145,22uit133
// 	dept := fmt.Sprintf("%c"+"%c", roll_no[3], roll_no[4]) // find the specific table for that roll_no
// 	dept_students_table := dept + "_students"
// 	dept_chatlist_table := dept + "_chatlist"
// 	return dept, dept_students_table, dept_chatlist_table
// }

// func FindDeptChatlistByRollno(roll_no string) string {
// 	dept := fmt.Sprintf("%c"+"%c", roll_no[3], roll_no[4]) // find the specific table for that roll_no
// 	dept_table := dept + "_chatlist"
// 	return dept_table
// }

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
	email = strings.TrimSpace(email)
	// Student regex
	studentRe := regexp.MustCompile(`^(\d{2})([a-z]+)(\d+)@kamarajengg\.edu\.in$`)
	// Staff regex
	staffRe := regexp.MustCompile(`^([a-z]+)(cse|ece|it|eee|mech|mtr|civil|ads|bt)@kamarajengg\.edu\.in$`)
	var whose_email string
	if studentRe.MatchString(email) {
		whose_email = "Student"
	}

	if staffRe.MatchString(email) {
		whose_email = "Staff"
	}
	fmt.Println("whose email - ", whose_email)
	return whose_email
}
func Find_dept_from_staff_email(staff_email string) string {
	staffRe := regexp.MustCompile(`^([a-z]+)(cse|ece|it|eee|mech|mtr|civil|ads|bt)@kamarajengg\.edu\.in$`)
	mail_array := staffRe.FindStringSubmatch(staff_email)
	department := mail_array[2] // department is the 2nd group
	fmt.Printf("dept from email - %v", strings.ToUpper(department))
	return strings.ToUpper(department) //convert the lowercase dept to uppercase 'cse'->'CSE'
}

func Find_staff_or_student_by_id(ID string) string {
	Re := regexp.MustCompile(`^\d{2}[up](bt|cs|ec|it|ee|me|mt|ci|ad)\d{3}$`)
	if Re.MatchString(ID) {
		return "STUDENT"
	} else {
		return "STAFF"
	}

}

func Find_dept_from_groupId(group_id string) string {
	var dept string
	group_dept_regexp := regexp.MustCompile(`^(cs|ad|bt|ec|me|mt|it|ce|ee)_`)
	matches := group_dept_regexp.FindStringSubmatch(group_id) // 'cse_1' -> ['cs_','cs']
	if len(matches) > 1 {
		dept = matches[1] // ['cs_','cs'] -> 'cs'
	}
	return dept

}
