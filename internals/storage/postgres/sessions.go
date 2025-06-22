package postgres

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func GenerateSessionID() string {

	deleteQuery := "delete from Sessions where expires_at < now() "
	_, deleteErr := Db.Exec(deleteQuery)
	if deleteErr != nil {
		fmt.Println("Error while deleting the sesions")
	}

	sessionId := uuid.New().String()
	fmt.Println("sessionId - ", sessionId)
	query := "insert into sessions(session_id,user_id,expires_at) values($1,$2,$3)"
	_, err := Db.Exec(query, sessionId, 12, time.Now().Add(3*time.Hour))
	if err != nil {
		fmt.Println("session inserted failure ")
		return ""
	}
	return sessionId
}

func DeleteSession(session_id string) {
	query := "delete from Sessions where session_id = $1 "
	_, err := Db.Exec(query, session_id)
	if err != nil {
		fmt.Println("error while deleting session - ", err)
	}
}
