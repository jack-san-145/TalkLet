package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"tet/internals/models"
	"tet/internals/services"
	"tet/internals/storage/postgres"
	"time"

	"github.com/xuri/excelize/v2"
)

func GroupCreation(w http.ResponseWriter, r *http.Request) {

	// {"group_name":"xxx","group_members":[10,2,3,4]}

	isFound, AdminID := FindCookie(r)
	if !isFound {
		return
	}
	fmt.Println("AdminID - ", AdminID)
	var group models.NewGroup
	err := json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		fmt.Println("error while decode the group details - ", err)
		return
	}
	fmt.Printf("new group - %v and Admin - %v ", group, AdminID)
	group.Admin = append(group.Admin, AdminID)
	admin_byte, err := json.Marshal(group.Admin)
	if err != nil {
		fmt.Println("error while marshal the admin - ", err)
		return
	}

	postgres.NewGroupPDb(group, admin_byte)
}

func GroupCreationByExcel(w http.ResponseWriter, r *http.Request) {
	isFound, AdminID := FindCookie(r)
	if !isFound {
		return
	}
	isStaff := postgres.Verify_Staff(AdminID)
	if !isStaff {
		WriteJSON(w, r, map[string]string{"status": "Invalid staff "})
		return
	}
	var (
		file   multipart.File        // multipart is package and File is interface to do all the io operations
		header *multipart.FileHeader //FileHeader stores the Filename , size of the file
	)
	err := r.ParseMultipartForm(10 << 20) // for use 10mb of ram memory if the file size is more than that it will store the remaining part in the disk
	if err != nil {
		fmt.Println("error while parsing form - ", err)
		return
	}
	group_name := r.FormValue("group_name")
	fmt.Println("group_name - ", group_name)
	//FormFile used to parse form
	file, header, err = r.FormFile("excel_file")
	if err != nil {
		fmt.Println("error in the excel file - ", err)
		return
	}
	fileType := header.Header.Get("Content-Type")
	if fileType != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		http.Error(w, "Only .xlsx files are allowed", http.StatusBadRequest)
		return
	}
	fmt.Println("file - ", file)
	fmt.Println("file size - ", header.Size)
	fmt.Println("file Name - ", header.Filename)
	fmt.Println("admin - ", AdminID)
	SaveThisFile(&file, header) // function to save the execl file to disk
	StartCreateGroup(AdminID, group_name, &file, header)

}

func SaveThisFile(file *multipart.File, header *multipart.FileHeader) {
	//it creates the empty file
	destination, err := os.Create("../../Excel-files/" + header.Filename)
	if err != nil {
		fmt.Println("error while creating the destination - ", err)
		return
	}

	//To copy the original file contents into that empty file
	_, err = io.Copy(destination, *file)
	if err != nil {
		fmt.Println("error while copying excel file - ", err)
		return
	}
}

func Showfilecontents(file *multipart.File, header *multipart.FileHeader) {
	var Rows [][]string // rows is an 2D array with rows and columns
	excel_file, err := excelize.OpenFile("../../Excel-files/" + header.Filename)
	if err != nil {
		fmt.Println("error while open the file - ", err)
		return
	}
	fileName := excel_file.GetSheetName(0) //under the hood excel sheet are arranged as a list [0,1,2] -> 0 for "sheet1" ,...
	//access the rows by the sheet name , here GetRows returns the 2D array like [ [c1,c2,c3] ,[c1,c2,c3], [c1,c2,c3] ]
	Rows, err = excel_file.GetRows(fileName)
	fmt.Println("Rows - ", Rows)
	if err != nil {
		fmt.Println("error while accessing the rows - ", err)
	} else if len(Rows) == 0 {
		fmt.Println("Empty Excel !! ")
	} else {

		for index, row := range Rows {
			fmt.Println("each row - ", row)
			if len(row) == 0 {
				continue
			}
			Sno, _ := strconv.Atoi(row[0])
			if Sno >= 1 {
				fmt.Printf("excel row %v - %v", index, row)
			}

		}
	}
	defer excel_file.Close()

}

func StartCreateGroup(admin string, group_name string, file *multipart.File, header *multipart.FileHeader) {
	var Rows [][]string
	var (
	// Department string
	// UploadedBy string
	)
	var Students []models.StudentDetails
	excel_file, err := excelize.OpenFile("../../Excel-files/" + header.Filename)
	if err != nil {
		fmt.Printf("error while opening the file '%v' - %v", header.Filename, err)
		return
	}
	sheet_name := excel_file.GetSheetName(0)
	Rows, err = excel_file.GetRows(sheet_name)
	if err != nil {
		fmt.Println("error while accessing the rows in the excel - ", err)
		return
	}
	isKcetTalklet := strings.TrimSpace(Rows[0][0])
	if isKcetTalklet != "KCET-TALKLET" {
		fmt.Println("Invalid file to create group")
		return
	}
	// fmt.Println("rows - ", Rows)
	// fmt.Println("before branch .. ", Rows[2][1])
	branch := strings.TrimSpace(Rows[2][1])
	dept := strings.TrimSpace(Rows[3][1])
	class := strings.TrimSpace(Rows[4][1])
	batch := strings.TrimSpace(Rows[5][1])
	chairperson := strings.TrimSpace(Rows[6][1])
	// UploadedBy = Rows[7][1]
	current_year, section := services.Find_Year_And_Section(strings.TrimSpace(class))
	batch_year, passing_year := services.FindBatch(strings.TrimSpace(batch))

	for index, row := range Rows {
		var new_Student models.StudentDetails
		if len(row) == 0 {
			continue
		}
		if index >= 10 {
			fmt.Println("")
			new_Student.Register_no = strings.TrimSpace(row[1])
			new_Student.Roll_no = strings.TrimSpace(row[2])
			new_Student.Name = strings.TrimSpace(row[3])
			new_Student.Email = strings.TrimSpace(row[5])
			new_Student.Mentor = strings.TrimSpace(row[6])
			new_Student.Branch = branch
			new_Student.Chairperson = chairperson
			new_Student.Current_year = current_year
			new_Student.Section = section
			new_Student.Batch_year = batch_year
			new_Student.Passing_year = passing_year
			dob_str := strings.TrimSpace(row[4])
			parsed, err := time.Parse("02-01-2006", dob_str)
			if err != nil {
				fmt.Println("error in dob parsing - ", err)
			}
			new_Student.DOB = parsed.Format("2006-01-02") // converting 15-01-2006 -> 2006-01-15
			// fmt.Println("register_no - ", new_Student.Register_no)
			// fmt.Println("roll_no - ", new_Student.Roll_no)
			// fmt.Println("name - ", new_Student.Name)
			// fmt.Println("dob - ", new_Student.DOB)
			// fmt.Println("email - ", new_Student.Email)
			// fmt.Println("mentor - ", new_Student.Mentor)
			// fmt.Println("branch - ", new_Student.Branch)
			// fmt.Println("Department - ", Department)
			// fmt.Printf("class - %v-%v ", new_Student.Current_year, new_Student.Section)
			// fmt.Printf("batch year - %v\npassing year - %v\nchairperson - %v\nuploadedby - %v", new_Student.Batch_year, new_Student.Passing_year, new_Student.Chairperson, UploadedBy)
			// fmt.Println("")
			Students = append(Students, new_Student)
		}

	}

	postgres.CreateNewGroupPDB(admin, group_name, dept, Students)

}
