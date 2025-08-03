package postgres

import (
	"context"
	"fmt"
	"tet/internals/models"
	"time"
)

func NewGroupPDb(group models.NewGroup, admin []byte) {
	var group_id int
	fmt.Println("new group to insert  - ", group)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "insert into allgroups(group_name,admin) values($1,$2::jsonb) returning group_id "
	// _, err := Db.Exec(query, group.Name, group.Admin)
	err := pool.QueryRow(ctx, query, group.Name, admin).Scan(&group_id)
	if err != nil {
		fmt.Println("error while inserting new group - ", err)
		return
	}
	fmt.Println("group_id - ", group_id)

}
