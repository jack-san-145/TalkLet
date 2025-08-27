package main

import (
	"fmt"
	"net/http"
	"os"
	"tet/internals/handlers"
	"tet/internals/storage/minio"
	"tet/internals/storage/postgres"
	"tet/internals/storage/redis"
	"tet/internals/websocket"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Println("error while loading .env file", err)
		return
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// 	router.Use(cors.Handler(cors.Options{
	//     AllowedOrigins:   []string{"http://localhost:3000", "http://127.0.0.1:3000"},
	//     AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	//     AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	//     ExposedHeaders:   []string{"Link"},
	//     AllowCredentials: true,
	//     MaxAge:           300,
	// }))

	router.Mount("/Static/", http.StripPrefix("/Static/", http.FileServer(http.Dir("../../ui/Static"))))
	router.Post("/update", handlers.UpdateHandler)

	router.Get("/", handlers.ServeLogin)

	router.Post("/talklet/send-otp-student", handlers.SendOtpHandler_for_students)
	router.Post("/talklet/send-otp-staff", handlers.SendOtpHandler_for_staffs)

	router.Post("/talklet/verify-otp-student", handlers.VerifyOTP_for_student_handler)
	router.Post("/talklet/verify-otp-staff", handlers.VerifyOTP_for_staff_handler)

	router.Post("/talklet/set-password", handlers.SetPassword)

	router.Get("/talklet/serve-register", handlers.ServeRegister)
	router.Post("/talklet/new-register", handlers.AccountRegisterHandler)
	router.Get("/talklet/serve-index", handlers.ServeIndex)
	router.Get("/talklet/profile/{id}", handlers.ProfileHandler)

	router.Post("/talklet/validate-student-login", handlers.StudentLoginValidationHandler)
	router.Post("/talklet/validate-staff-login", handlers.StaffLoginValidationHandler)

	router.Get("/ws", websocket.UpgradeToWebsocket)

	router.Post("/talklet/register-new-staff", handlers.StaffRegistration)
	router.Post("/talklet/creategroup-with-excel", handlers.GroupCreationByExcel)

	router.Get("/talklet/chat-history/{contact_id}", handlers.LoadChatMessages)
	router.Get("/talklet/get-all-chatlist", handlers.Chatlist)

	router.Post("/talklet/create-new-group", handlers.GroupCreation)

	router.Post("/talklet/chat/file-upload", handlers.ChatFileUploads)
	postgres.ConnectToDb()    //connect to postgres
	redis.CreateRedisClient() //create redis client
	minio.CreateMinioClient() //create minio client

	// cs,ad,bt,ec,it,me,mt

	// postgres.DropChatlistTable("cs")
	// postgres.DropChatlistTable("ad")
	// postgres.DropChatlistTable("bt")
	// postgres.DropChatlistTable("ec")
	// postgres.DropChatlistTable("me")
	// postgres.DropChatlistTable("mt")
	// postgres.DropChatlistTable("it")

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
