package account

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rusalexch/loyalty-group/internal/gophermart/internal/app"
)

type auth interface {
	CheckToken(ctx context.Context, authToken string) (app.User, error)
}

type createOrderFunc func(ctx context.Context, userID int, orderID string) error

type balance struct {
	Current  float64 `json:"current"`
	Withdraw float64 `json:"withdrawn"`
}

type withdrawRequest struct {
	OrderID string  `json:"order"`
	Amount  float64 `json:"sum"`
}

type Config struct {
	Pool        *pgxpool.Pool
	Auth        auth
	CreateOrder createOrderFunc
}

type accountModule struct {
	pool        *pgxpool.Pool
	auth        auth
	createOrder createOrderFunc
}

func New(conf Config) *accountModule {
	module := &accountModule{
		pool:        conf.Pool,
		auth:        conf.Auth,
		createOrder: conf.CreateOrder,
	}
	module.createTable()

	return module
}

func (am *accountModule) Middlewares() []app.Middleware {
	return []app.Middleware{}
}

func (am *accountModule) Handlers() []app.Handler {
	return []app.Handler{
		{
			Method:  http.MethodGet,
			Pattern: "/api/user/balance",
			Handler: am.balance,
		},
		{
			Method:  http.MethodPost,
			Pattern: "/api/user/balance/withdraw",
			Handler: am.withdraw,
		},
		{
			Method:  http.MethodGet,
			Pattern: "/api/user/withdrawals",
			Handler: am.withdrawals,
		},
	}
}

func (am *accountModule) Add(ctx context.Context, orderID string, amount float64) error {
	return am.addDebit(ctx, orderID, amount)
}
