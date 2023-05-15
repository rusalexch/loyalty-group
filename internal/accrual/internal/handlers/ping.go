package handlers

import (
	"log"
	"net/http"
)

func (h *handlers) ping(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := h.service.Ping(ctx); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
