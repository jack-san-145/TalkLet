package main

import (
	"fmt"
	"net/http"
	"os"
	"tet/internals/handlers"
	"tet/internals/storage/postgres"
	"tet/internals/storage/redis"
	"tet/internals/websocket"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Println("error while loading .env file", err)
		return
	}

	router := chi.NewRouter()
	router.Mount("/Static/", http.StripPrefix("/Static/", http.FileServer(http.Dir("../../ui/Static"))))
	router.Post("/update", handlers.UpdateHandler)
	router.Get("/", handlers.ServeLogin)
	router.Post("/talklet/send-otp-register", handlers.SendOtpRegisterHandler)
	router.Get("/talklet/serve-register", handlers.ServeRegister)
	router.Post("/talklet/new-register", handlers.AccountRegisterHandler)
	router.Post("/talklet/validate-login", handlers.LoginValidationHandler)
	router.Get("/talklet/serve-index", handlers.ServeIndex)
	router.Get("/talklet/profile/{id}", handlers.ProfileHandler)
	router.Get("/ws", websocket.UpgradeToWebsocket)

	postgres.ConnectToDb()    //connect to postgres
	redis.CreateRedisClient() //create redis client

	fmt.Println("server running")
	port := ":" + os.Getenv("PORT")
	err = http.ListenAndServe(port, router)
	if err != nil {
		fmt.Println("server failure - ", err)
	}
}
