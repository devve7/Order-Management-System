package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(h *Handler) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/orders", h.CreateOrder).Methods(http.MethodPost)

	return router
}
