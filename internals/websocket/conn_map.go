package websocket

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

var ConnMap sync.Map

func Set_ws_conn(id string, conn *websocket.Conn) {
	ConnMap.Store(id, conn) //stores the new websocket connection to the ConnMap

}

func Get_ws_Conn(id string) (*websocket.Conn, bool) {

	var (
		ws_conn  *websocket.Conn
		is_alive bool
		val      any
	)
	if val, is_alive = ConnMap.Load(id); is_alive {
		ws_conn = val.(*websocket.Conn)

	} else {
		fmt.Println("there is no websocket conection for this id - ")
	}
	return ws_conn, is_alive
}

func del_ws_conn(id string) {
	ConnMap.Delete(id) //deleting the web_socket connection form the ConnMap
}
