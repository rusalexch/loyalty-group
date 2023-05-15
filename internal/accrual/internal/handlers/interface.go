package handlers

import (
	"context"

	"github.com/rusalexch/loyalty-group/internal/accrual/internal/app"
)

type service interface {
	Ping(ctx context.Context) error
	GetOrder(ctx context.Context, orderID int64) (app.Order, error)
	AddReward(ctx context.Context, product app.Reward) error
	IsRewardExist(ctx context.Context, rewardID string) (bool, error)
	AddOrder(ctx context.Context, order app.OrderGoods) error
}
