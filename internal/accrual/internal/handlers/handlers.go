package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
)

type handlers struct {
	address string
	service service
	timeout time.Duration
	mux     *chi.Mux
}

func New(address string, service service) *handlers {
	h := &handlers{
		address: address,
		service: service,
		timeout: 5 * time.Second,
		mux:     chi.NewMux(),
	}

	return h
}

func (h *handlers) Start() {
	logger := httplog.NewLogger("httplog", httplog.Options{
		JSON: true,
	})
	h.mux.Use(middleware.RequestID)
	h.mux.Use(middleware.RealIP)
	h.mux.Use(httplog.RequestLogger(logger))
	h.mux.Use(middleware.Compress(5, "application/json", "text/html"))
	h.mux.Use(middleware.Recoverer)

	h.mux.Get("/", h.ping)
	h.mux.Get("/api/orders/{orderID}", h.calc)
	h.mux.Post("/api/orders", h.addOrder)
	h.mux.Post("/api/goods", h.addProduct)

	err := http.ListenAndServe(h.address, h.mux)
	if err != nil {
		log.Panic(err)
	}
}
