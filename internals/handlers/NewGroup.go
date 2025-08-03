package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"tet/internals/models"
	"tet/internals/storage/postgres"
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

func GroupCreationByExel(w http.ResponseWriter, r *http.Request) {
	isFound, AdminID := FindCookie(r)
	if !isFound {
		return
	}
	var (
		file   multipart.File
		header *multipart.FileHeader
	)
	err := r.ParseForm()
	if err != nil {
		fmt.Println("error while parsing form - ", err)
		return
	}
	file, header, err = r.FormFile("exel_file")
	if err != nil {
		fmt.Println("error in the exel file - ", err)
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
	SaveThisFile(&file, header)

}

func SaveThisFile(file *multipart.File, header *multipart.FileHeader) {
	destination, err := os.Create("../../Exel-files/" + header.Filename)
	if err != nil {
		fmt.Println("error while creating the destination - ", err)
		return
	}

	_, err = io.Copy(destination, *file)
	if err != nil {
		fmt.Println("error while copying exel file - ", err)
		return
	}
}
