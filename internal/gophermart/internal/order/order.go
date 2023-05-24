package order

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-chi/chi/v5"
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
}

type Config struct {
	Mux            *chi.Mux
	Pool           *pgxpool.Pool
	Auth           auth
	AccrualAddress string
}

type orderModule struct {
	mux            *chi.Mux
	pool           *pgxpool.Pool
	auth           auth
	account        balancer
	tick           *time.Ticker
	accrualAddress string
}

func New(conf Config) *orderModule {
	module := &orderModule{
		mux:            conf.Mux,
		pool:           conf.Pool,
		auth:           conf.Auth,
		tick:           time.NewTicker(2 * time.Second),
		accrualAddress: conf.AccrualAddress,
	}
	module.init()

	acc := account.New(account.Config{
		Mux:   conf.Mux,
		Pool:  conf.Pool,
		Auth:  conf.Auth,
		Order: module,
	})
	module.account = acc

	return module
}

func (om *orderModule) init() {
	om.initRepository()
	om.initHandler()
	go om.process()
}

func (om *orderModule) Create(ctx context.Context, userID int, orderID string) error {
	return om.add(ctx, userID, orderID)
}


