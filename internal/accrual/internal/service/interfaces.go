package service

import (
	"context"

	"github.com/rusalexch/loyalty-group/internal/accrual/internal/common"
)

type storager interface {
	Ping(ctx context.Context) error
}

type orderRepository interface {
	Add(ctx context.Context, orderID int64) error
	FindByID(ctx context.Context, orderID int64) (common.Order, error)
	UpdateStatus(ctx context.Context, orderID int64, status common.OrderStatus) error
	Update(ctx context.Context, order common.Order) error
	Delete(ctx context.Context, orderID int64) error
	FindRegistered(ctx context.Context) ([]int64, error)
}

type productRepository interface {
	Add(ctx context.Context, orderID int64, product []common.OrderProduct) error
	FindByOrderID(ctx context.Context, orderID int64) (common.OrderGoods, error)
}

type rewardRepository interface {
	Add(ctx context.Context, reward common.Reward) error
	Find(ctx context.Context, description string) (common.Reward, error)
}

type servicer interface {
	Ping(ctx context.Context) error
	GetOrder(ctx context.Context, orderID string) (float64, error)
	AddReward(ctx context.Context, product common.Reward) error
	AddOrder(ctx context.Context, order common.OrderGoods) error
}