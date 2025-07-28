package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, r *http.Request, data any) {
	data_byte, err := json.Marshal(data)
	if err != nil {
		fmt.Println("error while marshal data - ", err)
		return
	}
	w.Write(data_byte)
}
