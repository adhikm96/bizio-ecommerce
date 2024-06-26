package server

import (
	"github.com/Digital-AIR/bizio-ecommerce/internal/database"
	"github.com/Digital-AIR/bizio-ecommerce/internal/handler"
	"github.com/Digital-AIR/bizio-ecommerce/internal/handler/review"
	"log/slog"
	"net/http"
)

func InitServer() {
	database.MigrateDBSchema()
	StartServer()
}

func StartServer() {
	http.HandleFunc("/ping", handler.PingHandler)

	//Review handler
	http.HandleFunc("POST /api/v1/products/{product_id}/reviews", review.CreateReviewHanlder)
	http.HandleFunc("GET /api/v1/products/{product_id}/reviews", review.FetchReviewHandler)

	slog.Info("starting server at :8000")

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		slog.Error(err.Error())
	}
}
