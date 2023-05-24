package gophermart

import (
	"context"
	"log"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rusalexch/loyalty-group/internal/gophermart/internal/auth"
	"github.com/rusalexch/loyalty-group/internal/gophermart/internal/order"
	"github.com/rusalexch/loyalty-group/internal/gophermart/internal/server"
	"github.com/rusalexch/loyalty-group/internal/gophermart/internal/user"
)

type serv interface {
	Start() error
}

type Config struct {
	Address        string
	DBURL          string
	AccrualAddress string
}

type gophermart struct {
	s serv
}

func New(conf Config) *gophermart {
	mux := chi.NewMux()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	logger := httplog.NewLogger("httplog", httplog.Options{
		JSON: true,
	})
	mux.Use(httplog.RequestLogger(logger))
	mux.Use(middleware.Timeout(10 * time.Second))
	mux.Use(middleware.Compress(5, "application/json"))
	mux.Use(middleware.Recoverer)

	pool, err := pgxpool.New(context.Background(), conf.DBURL)
	if err != nil {
		log.Println("gophermart > can't create database pool")
		log.Fatal(err)
	}
	u := user.New(user.Config{
		Pool: pool,
	})
	a := auth.New(auth.Config{
		Mux:         mux,
		UserService: u,
		JwtSecret:   "super_secret",
	})
	order.New(order.Config{
		Mux:            mux,
		Pool:           pool,
		Auth:           a,
		AccrualAddress: conf.AccrualAddress,
	})

	s := server.New(conf.Address, mux)

	return &gophermart{
		s: s,
	}
}

func (g *gophermart) Start() {
	g.s.Start()
}
