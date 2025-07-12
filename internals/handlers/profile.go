package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"tet/internals/models"

	"github.com/go-chi/chi/v5"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	var profile models.Profile
	templ := template.Must(template.New("profile.html").Funcs(MyJsonConvFunc).ParseFiles("../../ui/templates/profile.html"))

	UserId_str := chi.URLParam(r, "id")
	fmt.Println("profile UserId - ", UserId_str)
	userId_int, _ := strconv.Atoi(UserId_str)
	profile.UserID = userId_int
	profile.Name = "jack"
	profile.Mobile = "7845941042"
	profile.Email = "jack145@gmail.com"
	profile.About = "my golden heart"
	profile.Joined = "Jan 2025"

	templ.Execute(w, profile)

}
