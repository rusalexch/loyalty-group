package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func (h *handlers) calc(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	orderID := chi.URLParam(r, "orderID")
	log.Println(ctx, orderID)
	// TODO

	w.WriteHeader(http.StatusOK)
}
