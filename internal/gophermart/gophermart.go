package gophermart

import (
	"context"
	"log"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rusalexch/loyalty-group/internal/gophermart/auth"
	"github.com/rusalexch/loyalty-group/internal/gophermart/order"
	"github.com/rusalexch/loyalty-group/internal/gophermart/server"
	"github.com/rusalexch/loyalty-group/internal/gophermart/user"
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
	})
	o := order.New(order.Config{
		Mux:            mux,
		Pool:           pool,
		Auth:           a,
		AccrualAddress: conf.AccrualAddress,
	})
	defer o.Close()

	s := server.New(conf.Address, mux)

	return &gophermart{
		s: s,
	}
}

func (g *gophermart) Start() {
	g.s.Start()
}
