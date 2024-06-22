package server

import (
	"github.com/Digital-AIR/bizio-ecommerce/internal/common"
	"github.com/Digital-AIR/bizio-ecommerce/internal/database"
	"github.com/Digital-AIR/bizio-ecommerce/internal/handler"
	"github.com/Digital-AIR/bizio-ecommerce/internal/session"
	"log"
	"net/http"
)

func InitServer() {
	database.MigrateDBSchema()
	StartServer()
}

func StartServer() {
	// user handler
	http.HandleFunc("POST /register", handler.HandleUserCreate)
	http.HandleFunc("POST /login", handler.HandleUserLogin)
	//http.HandleFunc("POST /logout", handler.HandleUserLogout)
	http.HandleFunc("/ping", pingHandler)

	log.Println("server started at :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func pingHandler(writer http.ResponseWriter, request *http.Request) {
	if !session.IsAuthenticated(request) {
		common.HandleUnAuthRes(writer, nil)
		return
	}

	writer.Write([]byte("{\"message\":\"pong\"}"))
}
