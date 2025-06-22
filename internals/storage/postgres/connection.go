package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

var Db *sql.DB

func ConnectToDb() {
	var err error
	conn := os.Getenv("POSTGRES_DATABASE_CONNECTION")
	fmt.Println("conn - ", conn)
	Db, err = sql.Open("postgres", conn)
	if err != nil {
		panic(err)
	}
	err = Db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("database connected successfully")
}
