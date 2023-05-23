package gophermart

import (
	"context"
	"log"

	"github.com/go-chi/chi/v5"
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
