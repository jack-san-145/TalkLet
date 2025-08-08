package postgres

import (
	"context"
	"fmt"
	"strings"
)

func SetPasswordDB(email string, password string) {
	roll_no := strings.Split(email, "@")[0]
	dept := fmt.Sprintf("%c"+"%c", email[3], email[4])
	dept_table := dept + "_students"
	fmt.Println("dept table - ", dept_table)
	fmt.Println("roll no - ", roll_no)
	query := fmt.Sprintf(`update %s set password = $1 where roll_no = $2 `, dept_table)
	_, err := pool.Exec(context.Background(), query, password, roll_no)
	if err != nil {
		fmt.Println("error while updating the password - ", err)
		return
	}
}
