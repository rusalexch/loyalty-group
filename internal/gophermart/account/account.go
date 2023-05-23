package account

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rusalexch/loyalty-group/internal/gophermart/app"
)

type auth interface {
	CheckToken(ctx context.Context, authToken string) (app.User, error)
}

type balance struct {
	Current  float64 `json:"current"`
	Withdraw float64 `json:"withdrawn"`
}

type withdrawRequest struct {
	OrderID string  `json:"order"`
	Amount  float64 `json:"sum"`
}

type Config struct {
	Mux  *chi.Mux
	Pool *pgxpool.Pool
	Auth auth
}

type accountModule struct {
	mux  *chi.Mux
	pool *pgxpool.Pool
	auth auth
}

func New(conf Config) *accountModule {
	module := &accountModule{
		mux:  conf.Mux,
		pool: conf.Pool,
		auth: conf.Auth,
	}
	module.init()

	return module
}

func (am *accountModule) init() {
	am.createTable()
	am.initHandler()
}

func (am *accountModule) Add(ctx context.Context, orderID string, amount float64) error {
	return am.addDebit(ctx, orderID, amount)
}
