package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(h *Handler) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/orders", h.CreateOrder).Methods(http.MethodPost)
	router.HandleFunc("/orders/{id:[0-9]+}", h.GetOrder).Methods(http.MethodGet)
	router.HandleFunc("/orders", h.GetOrders).Methods(http.MethodGet)

	router.HandleFunc("/orders/{id:[0-9]+}/items", h.AddItem).Methods(http.MethodPost)
	router.HandleFunc("/orders/{id:[0-9]+}/items/{item_id:[0-9]+}", h.DeleteItem).Methods(http.MethodDelete)

	router.HandleFunc("/orders/{id:[0-9]+}/pay", h.PayOrder).Methods(http.MethodPatch)
	router.HandleFunc("/orders/{id:[0-9]+}/ship", h.ShipOrder).Methods(http.MethodPatch)
	router.HandleFunc("/orders/{id:[0-9]+}/cancel", h.CancelOrder).Methods(http.MethodPatch)

	return router
}
