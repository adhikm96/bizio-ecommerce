package server

import (
	"github.com/Digital-AIR/bizio-ecommerce/internal/database"
	"github.com/Digital-AIR/bizio-ecommerce/internal/handler"
	"log/slog"
	"net/http"
)

func InitServer() {
	database.MigrateDBSchema()
	StartServer()
}

func StartServer() {
	http.HandleFunc("/ping", handler.PingHandler)

	slog.Info("starting server at :8000")

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		slog.Error(err.Error())
	}
}
