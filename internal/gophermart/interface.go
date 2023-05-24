package gophermart

import "github.com/rusalexch/loyalty-group/internal/gophermart/internal/app"

type moduler interface {
	Middlewares() []app.Middleware
	Handlers() []app.Handler
}
