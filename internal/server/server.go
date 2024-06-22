package server

import (
	"github.com/Digital-AIR/bizio-ecommerce/internal/database"
	"github.com/Digital-AIR/bizio-ecommerce/internal/handler"
	"log"
	"net/http"
)

func InitServer() {
	database.MigrateDBSchema()
	StartServer()
}

func StartServer() {
	http.HandleFunc("/ping", handler.PingHandler)

	log.Println("server started at :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
