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
	router.Post("/talklet/send-otp", handlers.SendOtpHandler)
	router.Post("/talklet/verify-otp", handlers.VerifyOTP)
	router.Post("/talklet/set-password", handlers.SetPassword)
	router.Get("/talklet/serve-register", handlers.ServeRegister)
	router.Post("/talklet/new-register", handlers.AccountRegisterHandler)
	router.Post("/talklet/validate-login", handlers.LoginValidationHandler)
	router.Get("/talklet/serve-index", handlers.ServeIndex)
	router.Get("/talklet/profile/{id}", handlers.ProfileHandler)
	router.Get("/ws", websocket.UpgradeToWebsocket)

	// router.Get("/talklet/chat-history/{contact_id}", handlers.LoadChatMessages)
	// router.Get("/talklet/get-all-OTO-chatlist", handlers.OneToOneChatlist)

	router.Post("/talklet/create-new-group", handlers.GroupCreation)
	router.Post("/talklet/creategroup-with-excel", handlers.GroupCreationByExcel)

	postgres.ConnectToDb()    //connect to postgres
	redis.CreateRedisClient() //create redis client

	// cs,ad,bt,ec,it,me,mt
	// postgres.DropAllTable()
	// go postgres.AddNewDepartment("cs")
	// go postgres.AddNewDepartment("ad")
	// go postgres.AddNewDepartment("bt")
	// go postgres.AddNewDepartment("ec")
	// go postgres.AddNewDepartment("me")
	// go postgres.AddNewDepartment("mt")
	// go postgres.AddNewDepartment("it")
	// go postgres.AddNewDepartment("civil")
	// go postgres.AddNewDepartment("eee")

	// postgres.AlterTable()
	fmt.Println("server running")
	port := ":" + os.Getenv("PORT")
	err = http.ListenAndServe(port, router)
	if err != nil {
		fmt.Println("server failure - ", err)
	}
}
