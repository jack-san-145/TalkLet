package postgres

import (
	"fmt"
	"tet/internals/models"
)

func NewGroupPDb(group models.NewGroup, admin []byte) {
	var group_id int
	fmt.Println("new group to insert  - ", group)
	query := "insert into allgroups(group_name,admin) values($1,$2::jsonb) returning group_id "
	// _, err := Db.Exec(query, group.Name, group.Admin)
	err := Db.QueryRow(query, group.Name, admin).Scan(&group_id)
	if err != nil {
		fmt.Println("error while inserting new group - ", err)
		return
	}
	fmt.Println("group_id - ", group_id)
}
