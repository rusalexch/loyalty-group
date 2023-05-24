package order

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rusalexch/loyalty-group/internal/gophermart/internal/account"
	"github.com/rusalexch/loyalty-group/internal/gophermart/internal/app"
)

type accrual struct {
	ID      string   `json:"order"`
	Status  string   `json:"status"`
	Accrual *float64 `json:"accrual,omitempty"`
}

type order struct {
	ID         string          `db:"id"`
	UserID     int             `db:"user_id"`
	Status     string          `db:"status"`
	Accrual    sql.NullFloat64 `db:"accrual"`
	UploadedAt time.Time       `db:"uploaded_at"`
}

type updateOrder struct {
	ID      string
	Status  string
	Accrual sql.NullFloat64
}

type auth interface {
	CheckToken(ctx context.Context, authToken string) (app.User, error)
}

type balancer interface {
	Add(ctx context.Context, orderID string, amount float64) error
	Middlewares() []app.Middleware
	Handlers() []app.Handler
}

type Config struct {
	Pool           *pgxpool.Pool
	Auth           auth
	AccrualAddress string
}

type orderModule struct {
	pool           *pgxpool.Pool
	auth           auth
	account        balancer
	tick           *time.Ticker
	accrualAddress string
}

func New(conf Config) *orderModule {
	module := &orderModule{
		pool:           conf.Pool,
		auth:           conf.Auth,
		tick:           time.NewTicker(2 * time.Second),
		accrualAddress: conf.AccrualAddress,
	}
	module.initRepository()

	acc := account.New(account.Config{
		Pool:        conf.Pool,
		Auth:        conf.Auth,
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
