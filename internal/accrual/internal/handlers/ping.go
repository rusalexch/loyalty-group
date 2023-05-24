package handlers

import (
	"log"
	"net/http"
)

func (h *handlers) ping(w http.ResponseWriter, r *http.Request) {
	log.Println("ping1")
	ctx := r.Context()
	log.Println("ping2")
	log.Println(ctx)

	if err := h.service.Ping(ctx); err != nil {
		log.Println("handlers > ping > can't ping")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte("pong"))
}
