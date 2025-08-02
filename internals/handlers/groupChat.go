package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
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
