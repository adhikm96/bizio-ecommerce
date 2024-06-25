package server

import (
	"github.com/Digital-AIR/bizio-ecommerce/internal/database"
	"github.com/Digital-AIR/bizio-ecommerce/internal/handler"
	"github.com/Digital-AIR/bizio-ecommerce/internal/handler/notification"
	"log/slog"
	"net/http"
)

func InitServer() {
	database.MigrateDBSchema()
	StartServer()
}

func StartServer() {

	http.HandleFunc("/api/v1/ping", handler.PingHandler)

	// notification api
	http.Handle("POST /api/v1/notifications/{user_id}", JSONHeaderMiddleware(http.HandlerFunc(notification.CreateHandler)))
	http.Handle("GET /api/v1/notifications/{user_id}", JSONHeaderMiddleware(http.HandlerFunc(notification.UsersNotificationHandler)))
	http.Handle("GET /api/v1/notifications/{id}/read", JSONHeaderMiddleware(http.HandlerFunc(notification.GetHandler)))

	slog.Info("starting server at :8000")

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		slog.Error(err.Error())
	}
}
