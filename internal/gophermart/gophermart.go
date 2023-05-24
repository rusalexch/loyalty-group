package gophermart

import (
	"context"
	"log"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rusalexch/loyalty-group/internal/gophermart/internal/app"
	"github.com/rusalexch/loyalty-group/internal/gophermart/internal/auth"
	"github.com/rusalexch/loyalty-group/internal/gophermart/internal/order"
	"github.com/rusalexch/loyalty-group/internal/gophermart/internal/server"
	"github.com/rusalexch/loyalty-group/internal/gophermart/internal/user"
)

type moduler interface {
	Middlewares() []app.Middleware
	Handlers() []app.Handler
}

type serv interface {
	Start() error
}

type Config struct {
	Address        string
	DBURL          string
	AccrualAddress string
}

type gophermart struct {
	s    serv
	mux  *chi.Mux
	mid  []app.Middleware
	hand []app.Handler
}

func New(conf Config) *gophermart {
	g := &gophermart{
		mid:  make([]app.Middleware, 0),
		hand: make([]app.Handler, 0),
		mux:  chi.NewMux(),
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
		Auth:           a,
		AccrualAddress: conf.AccrualAddress,
	})
	g.use(o)

	g.initRoute()

	g.s = server.New(conf.Address, g.mux)

	return g
}

func (g *gophermart) Start() {
	g.s.Start()
}

func (g *gophermart) use(module moduler) {
	mid := module.Middlewares()
	hand := module.Handlers()
	g.mid = append(g.mid, mid...)
	g.hand = append(g.hand, hand...)
}

func (g *gophermart) addAppMiddleware() {
	for _, m := range g.mid {
		g.mux.Use(m)
	}
}

func (g *gophermart) addAppHandler() {
	for _, h := range g.hand {
		g.mux.MethodFunc(h.Method, h.Pattern, h.Handler)
	}
}

func (g *gophermart) initRoute() {
	g.mux.Use(middleware.RequestID)
	g.mux.Use(middleware.RealIP)
	logger := httplog.NewLogger("httplog", httplog.Options{
		JSON: true,
	})
	g.mux.Use(httplog.RequestLogger(logger))
	g.mux.Use(middleware.Timeout(10 * time.Second))
	g.addAppMiddleware()
	g.mux.Use(middleware.Compress(5, "application/json"))
	g.mux.Use(middleware.Recoverer)

	g.addAppHandler()
}
