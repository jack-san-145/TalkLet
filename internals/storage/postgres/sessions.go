package postgres

import (
	"fmt"
	"tet/internals/models"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func GenerateSessionID(userId int) models.Session {
	var session models.Session
	deleteQuery := "delete from Sessions where expires_at < now() "
	_, deleteErr := Db.Exec(deleteQuery)
	if deleteErr != nil {
		fmt.Println("Error while deleting the sesions")
	}

	sessionId := uuid.New().String()
	fmt.Println("sessionId - ", sessionId)
	query := "insert into sessions(session_id,user_id,expires_at) values($1,$2,$3) returning * ;"
	err := Db.QueryRow(query, sessionId, userId, time.Now().Add(3*time.Hour)).Scan(&session.Session_id, &session.User_id, &session.Expires_at)
	if err != nil {
		fmt.Println("session inserted failure ", err)
	}
	return session
}

func FindSessionPdb(session_id string) (int, models.Session) {
	var session models.Session
	query := "select * from Sessions where session_id = $1 "
	err := Db.QueryRow(query, session_id).Scan(&session.Session_id, &session.User_id, &session.Expires_at)
	if err != nil {
		fmt.Println("error while find session_id in postgers - ", err)
	}
	return session.User_id, session
}

func DeleteSession(session_id string) {
	query := "delete from Sessions where session_id = $1 "
	_, err := Db.Exec(query, session_id)
	if err != nil {
		fmt.Println("error while deleting session - ", err)
	}
}
