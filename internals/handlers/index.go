package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("../../ui/templates/index.html")
	if err != nil {
		fmt.Println("file not found - ", err)
		return
	}
	err = temp.Execute(w, nil)
	if err != nil {
		fmt.Println("error while executing - ", err)
	}
}
