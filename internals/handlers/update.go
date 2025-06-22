package handlers

import (
	"fmt"
	"net/http"
)

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("error while parsing form")
		return
	}
	data := r.FormValue("msg")
	fmt.Fprintf(w, "<h3>%s</h3>", data)
}
