package handlers

import (
	"context"

	"github.com/rusalexch/loyalty-group/internal/accrual/internal/common"
)

type service interface {
	Ping(ctx context.Context) error
	GetOrder(ctx context.Context, orderID string) (float64, error)
	AddProduct(ctx context.Context, product common.Reward) error
	AddOrder(ctx context.Context, order common.OrderGoods) error
}
