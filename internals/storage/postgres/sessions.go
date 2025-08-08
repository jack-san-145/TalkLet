package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"tet/internals/models"
	"time"
)

func GenerateSessionID(roll_no string) models.Session {
	var session models.Session

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	deleteQuery := "delete from Sessions where expires_at < now() "
	_, deleteErr := pool.Exec(ctx, deleteQuery)
	if deleteErr != nil {
		fmt.Println("Error while deleting the sesions")
	}

	sessionId := uuid.New().String()
	fmt.Println("sessionId - ", sessionId)
	query := "insert into sessions(session_id,roll_no,expires_at) values($1,$2,$3) returning * ;"
	err := pool.QueryRow(ctx, query, sessionId, roll_no, time.Now().Add(3*time.Hour)).Scan(&session.Session_id, &session.Roll_no, &session.Expires_at)
	if err != nil {
		fmt.Println("session inserted failure ", err)
	}
	return session
}

func FindSessionPdb(session_id string) (string, models.Session) {
	var session models.Session

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "select * from Sessions where session_id = $1 "
	err := pool.QueryRow(ctx, query, session_id).Scan(&session.Session_id, &session.Roll_no, &session.Expires_at)
	if err != nil {
		fmt.Println("error while find session_id in postgers - ", err)
	}
	return session.Roll_no, session
}

func DeleteSession(session_id string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "delete from Sessions where session_id = $1 "
	_, err := pool.Exec(ctx, query, session_id)
	if err != nil {
		fmt.Println("error while deleting session - ", err)
	}
}
