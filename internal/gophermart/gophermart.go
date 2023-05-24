package gophermart

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rusalexch/loyalty-group/internal/gophermart/internal/app"
	"github.com/rusalexch/loyalty-group/internal/gophermart/internal/auth"
	"github.com/rusalexch/loyalty-group/internal/gophermart/internal/order"
	"github.com/rusalexch/loyalty-group/internal/gophermart/internal/user"
)

func New(conf Config) *gophermart {
	g := &gophermart{
		mid:  make([]app.Middleware, 0),
		hand: make([]app.Handler, 0),
		mux:  chi.NewMux(),
		addr: conf.Address,
	}

	pool, err := pgxpool.New(context.Background(), conf.DBURL)
	if err != nil {
		log.Println("gophermart > can't create database pool")
		log.Fatal(err)
	}
	u := user.New(user.Config{
		Pool: pool,
	})
	a := auth.New(auth.Config{
		UserService: u,
		JwtSecret:   "super_secret",
	})
	g.use(a)
	o := order.New(order.Config{
		Pool:           pool,
		AccrualAddress: conf.AccrualAddress,
	})
	g.use(o)

	g.initRoute()

	return g
}

func (g *gophermart) Start() {
	log.Println(http.ListenAndServe(g.addr, g.mux))
}
