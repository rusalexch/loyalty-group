package order

import (
	"context"

	"github.com/rusalexch/loyalty-group/internal/gophermart/internal/app"
)

type balancer interface {
	Add(ctx context.Context, orderID string, amount float64) error
	Middlewares() []app.Middleware
	Handlers() []app.Handler
}
