package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"tet/internals/models"
	"tet/internals/storage/postgres"

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
	var (
		file   multipart.File        // multipart is package and File is interface to do all the io operations
		header *multipart.FileHeader //FileHeader stores the Filename , size of the file
	)
	err := r.ParseForm()
	if err != nil {
		fmt.Println("error while parsing form - ", err)
		return
	}
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
	Showfilecontents(file, header)
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
	if err != nil {
		fmt.Println("error while accessing the rows - ", err)
	} else if len(Rows) == 0 {
		fmt.Println("Empty Excel !! ")
	} else {

		for index, row := range Rows {
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
