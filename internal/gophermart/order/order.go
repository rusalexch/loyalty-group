package order

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rusalexch/loyalty-group/internal/gophermart/account"
	"github.com/rusalexch/loyalty-group/internal/gophermart/app"
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := om.pool.Exec(ctx, sqlCreateTable)
	if err != nil {
		log.Println("ordfer > init > can't create user_orders table")
		log.Fatal(err)
	}

	om.mux.Get("/api/user/orders", om.get)
	om.mux.Post("/api/user/orders", om.create)

	go om.process()
}

func (om *orderModule) Create(ctx context.Context, userID int, orderID string) error {
	return om.add(ctx, userID, orderID)
}

func (om *orderModule) findByID(ctx context.Context, orderID string) (order, error) {
	var ord order
	row := om.pool.QueryRow(ctx, sqlFindByID, orderID)
	err := row.Scan(&ord.ID, &ord.UserID, &ord.Status, &ord.Accrual, &ord.UploadedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return order{}, app.ErrNotFound
		}
		return order{}, err
	}

	return ord, nil
}

func (om *orderModule) findByUserID(ctx context.Context, userID int) ([]order, error) {
	orders := make([]order, 0, 10)
	rows, err := om.pool.Query(ctx, sqlFundByUserID, userID)
	if err != nil {
		return []order{}, err
	}

	for rows.Next() {
		var ord order
		err := rows.Scan(&ord.ID, &ord.UserID, &ord.Status, &ord.Accrual, &ord.UploadedAt)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []order{}, app.ErrNotFound
			}
			return []order{}, err
		}
		orders = append(orders, ord)
	}

	return orders, nil
}

func (om *orderModule) add(ctx context.Context, userID int, orderID string) error {
	_, err := om.pool.Exec(ctx, sqlAdd, orderID, userID)

	return err
}

func (om *orderModule) findRegistered(ctx context.Context) ([]order, error) {
	orders := make([]order, 0, 10)
	rows, err := om.pool.Query(ctx, sqlFindRegistered)
	if err != nil {
		return []order{}, err
	}
	for rows.Next() {
		var ord order
		err := rows.Scan(&ord.ID, &ord.UserID, &ord.Status, &ord.Accrual, &ord.UploadedAt)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []order{}, app.ErrNotFound
			}
			return []order{}, err
		}
		orders = append(orders, ord)
	}

	return orders, nil
}

func (om *orderModule) updateOrder(ctx context.Context, updOrder updateOrder) error {
	_, err := om.pool.Exec(ctx, sqlUpdateOrder, updOrder.ID, updOrder.Status, updOrder.Accrual)

	return err
}

func dbToJSON(ord order) app.Order {
	var accrual *float64 = nil
	if ord.Accrual.Valid {
		accrual = &ord.Accrual.Float64
	}
	return app.Order{
		ID:         ord.ID,
		Status:     ord.Status,
		Accrual:    accrual,
		UploadedAt: ord.UploadedAt,
	}
}
