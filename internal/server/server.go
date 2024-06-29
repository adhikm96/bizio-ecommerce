package server

import (
	"github.com/Digital-AIR/bizio-ecommerce/internal/database"
	"github.com/Digital-AIR/bizio-ecommerce/internal/handler"
	"github.com/Digital-AIR/bizio-ecommerce/internal/handler/inventory"
	"github.com/Digital-AIR/bizio-ecommerce/internal/handler/notification"
	"github.com/Digital-AIR/bizio-ecommerce/internal/handler/order"
	"github.com/Digital-AIR/bizio-ecommerce/internal/handler/review"
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
	http.Handle("POST /api/v1/notifications", JSONHeaderMiddleware(http.HandlerFunc(notification.CreateHandler)))
	http.Handle("GET /api/v1/notifications/{user_id}", JSONHeaderMiddleware(http.HandlerFunc(notification.UsersNotificationHandler)))
	http.Handle("PUT /api/v1/notifications/{id}/read", JSONHeaderMiddleware(http.HandlerFunc(notification.UpdateReadNotification)))

	//Review api
	http.Handle("POST /api/v1/products/{product_id}/reviews", JSONHeaderMiddleware(http.HandlerFunc(review.CreateReviewHanlder)))
	http.Handle("GET /api/v1/products/{product_id}/reviews", JSONHeaderMiddleware(http.HandlerFunc(review.FetchReviewHandler)))

	// order api
	http.Handle("GET /api/v1/orders/{order_id}", JSONHeaderMiddleware(http.HandlerFunc(order.FetchOrderHandler)))
	http.Handle("POST /api/v1/orders/{order_id}", JSONHeaderMiddleware(http.HandlerFunc(order.CreateOrderHandler)))
	http.Handle("PUT /api/v1/orders/{order_id}", JSONHeaderMiddleware(http.HandlerFunc(order.UpdateOrderStatusHandler)))

	http.Handle("GET /api/v1/inventory/{variantId}", JSONHeaderMiddleware(http.HandlerFunc(inventory.FetchInventoryHandler)))
	http.Handle("PUT /api/v1/admin/inventory/{variantId}", JSONHeaderMiddleware(http.HandlerFunc(inventory.UpdateInventoryHandler)))

	slog.Info("starting server at :8000")

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		slog.Error(err.Error())
	}
}
