package service

import (
	"context"

	"github.com/rusalexch/loyalty-group/internal/accrual/internal/app"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go

type storager interface {
	Ping(ctx context.Context) error
}

type orderRepository interface {
	Add(ctx context.Context, orderID string) error
	FindByID(ctx context.Context, orderID string) (app.Order, error)
	UpdateStatus(ctx context.Context, orderID string, status string) error
	Update(ctx context.Context, order app.Order) error
	Delete(ctx context.Context, orderID string) error
	FindRegistered(ctx context.Context) ([]string, error)
}

type productRepository interface {
	Add(ctx context.Context, orderID string, product []app.OrderProduct) error
	FindByOrderID(ctx context.Context, orderID string) (app.OrderGoods, error)
}

type rewardRepository interface {
	Add(ctx context.Context, reward app.Reward) error
	Find(ctx context.Context, description string) (app.Reward, error)
	FindByID(ctx context.Context, ID string) (app.Reward, error)
}
