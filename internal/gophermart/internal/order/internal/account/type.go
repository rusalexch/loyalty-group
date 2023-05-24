package account

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

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
	CreateOrder createOrderFunc
}

type accountModule struct {
	pool        *pgxpool.Pool
	createOrder createOrderFunc
}
