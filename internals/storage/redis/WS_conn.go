package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

func Set_Conn_to_Redis(id string, conn *websocket.Conn) {
	err := rdb.Set(context.Background(), id, conn, time.Hour*24).Err()
	if err != nil {
		fmt.Println("error while setting the each connection to redis - ", err)
		return
	}
}

func Del_conn_from_Redis(id string) {
	err := rdb.Del(context.Background(), id).Err()
	if err != nil {
		fmt.Println("error while deleting the connection from the redis - ", err)
		return
	}

}

func Get_conn_from_Redis(id string) {
	err := rdb.Get(context.Background(), id).Err()
	if err != nil {
		fmt.Println("error while accessing the Conn from the redis - ", err)
		return
	}
}
