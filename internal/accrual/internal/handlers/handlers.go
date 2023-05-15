package handlers

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
)

type handlers struct {
	address      string
	service      service
	timeout      time.Duration
	mux          *chi.Mux
	maxRequest   int
	mx           *sync.Mutex
	reqPerMinute map[string]int
	ticker       *time.Ticker
}

func New(address string, service service) *handlers {
	h := &handlers{
		address: address,
		service: service,
		timeout: 5 * time.Second,
		mux:     chi.NewMux(),
		ticker:  time.NewTicker(time.Minute),
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
	h.mux.Use(middleware.Compress(5, "application/json"))
	h.mux.Use(middleware.Recoverer)

	h.mux.Get("/", h.ping)
	h.mux.Get("/api/orders/{orderID}", h.getOrder)
	h.mux.Post("/api/orders", h.addOrder)
	h.mux.Post("/api/goods", h.addReward)

	go func() {
		for range h.ticker.C {
			h.resetRequestCounter()
		}
	}()

	log.Println(http.ListenAndServe(h.address, h.mux))
}

func (h *handlers) resetRequestCounter() {
	h.mx.Lock()
	defer h.mx.Unlock()

	h.reqPerMinute = make(map[string]int)
}
