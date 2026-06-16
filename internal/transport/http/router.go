// Package http ...
package http

import (
	"Order-Management-System/internal/domain/user"
	"Order-Management-System/internal/transport/http/auth"
	health "Order-Management-System/internal/transport/http/health"
	"Order-Management-System/internal/transport/http/middleware"
	http_order "Order-Management-System/internal/transport/http/order"
	http_product "Order-Management-System/internal/transport/http/product"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func NewRouter(
	orderHandler *http_order.OrderHandler,
	productHandler *http_product.ProductHandler,
	healthHandler *health.HealthHandler,
	authHandler *auth.AuthHandler,
	authMiddleware *middleware.AuthMiddleware,
	logger *logrus.Logger,
) http.Handler {
	requireAuth := func(h http.HandlerFunc) http.Handler {
		return authMiddleware.RequireAuth(http.HandlerFunc(h))
	}
	requireAdmin := func(h http.HandlerFunc) http.Handler {
		return authMiddleware.RequireAuth(
			middleware.RequireRoleMiddleware(user.RoleAdmin, http.HandlerFunc(h)),
		)
	}

	router := mux.NewRouter()

	router.Use(middleware.LoggingMiddleWare(logger))

	// Health
	router.HandleFunc("/health", healthHandler.Health).Methods(http.MethodGet)

	//Orders
	router.Handle("/orders", requireAuth(orderHandler.CreateOrder)).Methods(http.MethodPost)

	router.Handle("/orders/{id:[0-9]+}", requireAuth(orderHandler.GetOrder)).Methods(http.MethodGet)
	router.Handle("/orders", requireAuth(orderHandler.GetOrders)).Methods(http.MethodGet)

	router.Handle("/orders/{id:[0-9]+}/items", requireAuth(orderHandler.AddItem)).Methods(http.MethodPost)
	router.Handle("/orders/{id:[0-9]+}/items/{item_id:[0-9]+}", requireAuth(orderHandler.DeleteItem)).Methods(http.MethodDelete)
	router.Handle("/orders/{id:[0-9]+}/pay", requireAuth(orderHandler.PayOrder)).Methods(http.MethodPost)
	router.Handle("/orders/{id:[0-9]+}/cancel", requireAuth(orderHandler.CancelOrder)).Methods(http.MethodPost)
	router.Handle("/orders/{id:[0-9]+}/ship", requireAdmin(orderHandler.ShipOrder)).Methods(http.MethodPost)

	// Products
	router.Handle("/products", requireAdmin(productHandler.CreateProduct)).Methods(http.MethodPost)

	router.HandleFunc("/products/{id:[0-9]+}", productHandler.GetProduct).Methods(http.MethodGet)
	router.HandleFunc("/products", productHandler.GetProducts).Methods(http.MethodGet)

	router.Handle("/products/{id:[0-9]+}/price", requireAdmin(productHandler.ChangePrice)).Methods(http.MethodPost)
	router.Handle("/products/{id:[0-9]+}/activate", requireAdmin(productHandler.ActivateProduct)).Methods(http.MethodPost)
	router.Handle("/products/{id:[0-9]+}/deactivate", requireAdmin(productHandler.DeactivateProduct)).Methods(http.MethodPost)
	router.Handle("/products/{id:[0-9]+}/stock/add", requireAdmin(productHandler.AddStock)).Methods(http.MethodPost)
	router.Handle("/products/{id:[0-9]+}/stock/remove", requireAdmin(productHandler.RemoveStock)).Methods(http.MethodPost)

	// Auth
	router.HandleFunc("/auth/register", authHandler.Register).Methods(http.MethodPost)
	router.HandleFunc("/auth/login", authHandler.Login).Methods(http.MethodPost)

	// Me
	router.Handle("/me", requireAuth(authHandler.Me)).Methods(http.MethodGet)

	return router
}
