package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"tet/internals/models"
	"tet/internals/storage/postgres"
)

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	isFound, userId := FindCookie(r)
	if !isFound {
		return
	}
	Chatlist := postgres.LoadChatlist(userId)
	templ := template.Must(template.New("index.html").Funcs(MyJsonConvFunc).ParseFiles("../../ui/templates/index.html"))
	// temp, err := template.ParseFiles("../../ui/templates/index.html")
	// if err != nil {
	// 	fmt.Println("file not found - ", err)
	// 	return
	// }
	fmt.Println(Chatlist)
	err := templ.Execute(w, map[string][]models.Chatlist{
		"Chat_list_go": Chatlist,
	})
	if err != nil {
		fmt.Println("error while executing - ", err)
	}
}
