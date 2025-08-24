package postgres

import (
	"context"
	"fmt"
	"strings"
	"tet/internals/services"
)

func SetPasswordDB(email string, password string) error {
	roll_no := strings.Split(email, "@")[0]
	_, dept_table, _ := services.Find_dept_from_rollNo(roll_no)

	fmt.Println("dept table - ", dept_table)
	fmt.Println("roll no - ", roll_no)
	query := fmt.Sprintf(`update %s set password = $1 where roll_no = $2 `, dept_table)
	_, err := pool.Exec(context.Background(), query, password, roll_no)
	if err != nil {
		fmt.Println("error while updating the password - ", err)
		return fmt.Errorf("err - %s", err)
	}
	return nil
}
