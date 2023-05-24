package account

import (
	"context"
	"net/http"

	"github.com/rusalexch/loyalty-group/internal/gophermart/internal/app"
)



func New(conf Config) *accountModule {
	module := &accountModule{
		pool:        conf.Pool,
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
