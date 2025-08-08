package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	// "tet/internals/storage/postgres"
)

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	isFound, roll_no := FindCookie(r)
	if !isFound {
		return
	}
	fmt.Println(roll_no)
	// Chatlist := postgres.LoadChatlist(userId)
	templ := template.Must(template.New("index.html").Funcs(MyJsonConvFunc).ParseFiles("../../ui/templates/index.html"))

	// fmt.Println(Chatlist)
	
	err := templ.Execute(w, nil)
	if err != nil {
		fmt.Println("error while executing - ", err)
	}
}
