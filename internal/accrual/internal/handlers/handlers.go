package handlers

import (
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
)

type handlers struct {
	service      service
	mux          *chi.Mux
	maxRequest   int
	mx           sync.Mutex
	reqPerMinute map[string]int
	ticker       *time.Ticker
}

func New(service service) http.Handler {
	h := &handlers{
		service:      service,
		mux:          chi.NewMux(),
		ticker:       time.NewTicker(time.Minute),
		reqPerMinute: map[string]int{},
		maxRequest:   60,
	}

	h.init()

	return h.mux
}

func (h *handlers) init() {
	logger := httplog.NewLogger("httplog", httplog.Options{
		JSON: true,
	})
	h.mux.Use(middleware.RequestID)
	h.mux.Use(middleware.RealIP)
	h.mux.Use(httplog.RequestLogger(logger))
	h.mux.Use(timeoutMiddleware)
	h.mux.Use(middleware.Compress(5, "application/json"))
	h.mux.Use(middleware.Recoverer)

	h.mux.Get("/ping", h.ping)
	h.mux.Get("/api/orders/{ID}", h.getOrder)
	h.mux.Post("/api/orders", h.addOrder)
	h.mux.Post("/api/goods", h.addReward)

	go func() {
		for range h.ticker.C {
			h.resetRequestCounter()
		}
	}()
}

func (h *handlers) resetRequestCounter() {
	h.mx.Lock()
	defer h.mx.Unlock()

	h.reqPerMinute = make(map[string]int)
}
