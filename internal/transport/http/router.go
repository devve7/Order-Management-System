// Package http ...
package http

import (
	health "Order-Management-System/internal/transport/http/health"
	"Order-Management-System/internal/transport/http/middleware"
	http_order "Order-Management-System/internal/transport/http/order"
	http_product "Order-Management-System/internal/transport/http/product"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func NewRouter(orderHandler *http_order.OrderHandler, productHandler *http_product.ProductHandler, healthHandler *health.HealthHandler, logger *logrus.Logger) http.Handler {
	router := mux.NewRouter()

	router.Use(middleware.LoggingMiddleWare(logger))

	router.HandleFunc("/health", healthHandler.Health).Methods(http.MethodGet)

	router.HandleFunc("/orders", orderHandler.CreateOrder).Methods(http.MethodPost)
	router.HandleFunc("/orders/{id:[0-9]+}", orderHandler.GetOrder).Methods(http.MethodGet)
	router.HandleFunc("/orders", orderHandler.GetOrders).Methods(http.MethodGet)

	router.HandleFunc("/orders/{id:[0-9]+}/items", orderHandler.AddItem).Methods(http.MethodPost)
	router.HandleFunc("/orders/{id:[0-9]+}/items/{item_id:[0-9]+}", orderHandler.DeleteItem).Methods(http.MethodDelete)

	router.HandleFunc("/orders/{id:[0-9]+}/pay", orderHandler.PayOrder).Methods(http.MethodPost)
	router.HandleFunc("/orders/{id:[0-9]+}/ship", orderHandler.ShipOrder).Methods(http.MethodPost)
	router.HandleFunc("/orders/{id:[0-9]+}/cancel", orderHandler.CancelOrder).Methods(http.MethodPost)

	router.HandleFunc("/products", productHandler.CreateProduct).Methods(http.MethodPost)
	router.HandleFunc("/products/{id:[0-9]+}", productHandler.GetProduct).Methods(http.MethodGet)
	router.HandleFunc("/products", productHandler.GetProducts).Methods(http.MethodGet)

	router.HandleFunc("/products/{id:[0-9]+}/price", productHandler.ChangePrice).Methods(http.MethodPost)
	router.HandleFunc("/products/{id:[0-9]+}/activate", productHandler.ActivateProduct).Methods(http.MethodPost)
	router.HandleFunc("/products/{id:[0-9]+}/deactivate", productHandler.DeactivateProduct).Methods(http.MethodPost)
	router.HandleFunc("/products/{id:[0-9]+}/stock/add", productHandler.AddStock).Methods(http.MethodPost)
	router.HandleFunc("/products/{id:[0-9]+}/stock/remove", productHandler.RemoveStock).Methods(http.MethodPost)

	return router
}
