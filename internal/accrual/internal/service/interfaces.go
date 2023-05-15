package service

import (
	"context"

	"github.com/rusalexch/loyalty-group/internal/accrual/internal/app"
)

type storager interface {
	Ping(ctx context.Context) error
}

type orderRepository interface {
	Add(ctx context.Context, orderID int64) error
	FindByID(ctx context.Context, orderID int64) (app.Order, error)
	UpdateStatus(ctx context.Context, orderID int64, status string) error
	Update(ctx context.Context, order app.Order) error
	Delete(ctx context.Context, orderID int64) error
	FindRegistered(ctx context.Context) ([]int64, error)
}

type productRepository interface {
	Add(ctx context.Context, orderID int64, product []app.OrderProduct) error
	FindByOrderID(ctx context.Context, orderID int64) (app.OrderGoods, error)
}

type rewardRepository interface {
	Add(ctx context.Context, reward app.Reward) error
	Find(ctx context.Context, description string) (app.Reward, error)
	FindByID(ctx context.Context, ID string) (app.Reward, error)
}
