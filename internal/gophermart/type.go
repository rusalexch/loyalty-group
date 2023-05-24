package gophermart

import (
	"github.com/go-chi/chi/v5"
	"github.com/rusalexch/loyalty-group/internal/gophermart/internal/app"
)

type Config struct {
	Address        string
	DBURL          string
	AccrualAddress string
}

type gophermart struct {
	addr string
	mux  *chi.Mux
	mid  []app.Middleware
	hand []app.Handler
}
