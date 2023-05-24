package gophermart

import (
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/httplog"
)

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
