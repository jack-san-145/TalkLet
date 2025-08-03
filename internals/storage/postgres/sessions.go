package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"tet/internals/models"
	"time"
)

func GenerateSessionID(userId int) models.Session {
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
	query := "insert into sessions(session_id,user_id,expires_at) values($1,$2,$3) returning * ;"
	err := pool.QueryRow(ctx, query, sessionId, userId, time.Now().Add(3*time.Hour)).Scan(&session.Session_id, &session.User_id, &session.Expires_at)
	if err != nil {
		fmt.Println("session inserted failure ", err)
	}
	return session
}

func FindSessionPdb(session_id string) (int, models.Session) {
	var session models.Session

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "select * from Sessions where session_id = $1 "
	err := pool.QueryRow(ctx, query, session_id).Scan(&session.Session_id, &session.User_id, &session.Expires_at)
	if err != nil {
		fmt.Println("error while find session_id in postgers - ", err)
	}
	return session.User_id, session
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
