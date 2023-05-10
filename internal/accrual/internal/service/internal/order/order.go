package order

import (
	"context"

	"github.com/rusalexch/loyalty-group/internal/accrual/internal/common"
)

type orderRepo interface {
	Add(ctx context.Context, orderID string) error
	FindById(ctx context.Context, orderID string) (common.Order, error)
	UpdateStatus(ctx context.Context, orderID string, status common.OrderStatus) error
	Update(ctx context.Context, order common.Order) error
}

type orderService struct {
	repo orderRepo
}

func New(repo orderRepo) *orderService {
	return &orderService{
		repo: repo,
	}
}
