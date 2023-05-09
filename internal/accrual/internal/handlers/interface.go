package handlers

import (
	"context"

	"github.com/rusalexch/loyalty-group/internal/accrual/internal/common"
)

type service interface {
	Add(ctx context.Context, product common.Reward) error
	Calc(ctx context.Context, products []common.OrderProduct) (float64, error)
	Ping(ctx context.Context) error
}
