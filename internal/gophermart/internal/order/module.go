package order

import (
	"net/http"
	"time"

	"github.com/rusalexch/loyalty-group/internal/gophermart/internal/app"
	"github.com/rusalexch/loyalty-group/internal/gophermart/internal/order/internal/account"
)

func New(conf Config) *orderModule {
	module := &orderModule{
		pool:           conf.Pool,
		tick:           time.NewTicker(2 * time.Second),
		accrualAddress: conf.AccrualAddress,
	}
	module.initRepository()

	acc := account.New(account.Config{
		Pool:        conf.Pool,
		CreateOrder: module.add,
	})
	module.account = acc

	go module.process()

	return module
}

func (om *orderModule) Middlewares() []app.Middleware {
	mid := om.account.Middlewares()

	return append([]app.Middleware{}, mid...)
}

func (om *orderModule) Handlers() []app.Handler {
	accHand := om.account.Handlers()
	hand := []app.Handler{
		{
			Method:  http.MethodGet,
			Pattern: "/api/user/orders",
			Handler: om.get,
		},
		{
			Method:  http.MethodPost,
			Pattern: "/api/user/orders",
			Handler: om.create,
		},
	}
	return append(hand, accHand...)
}
