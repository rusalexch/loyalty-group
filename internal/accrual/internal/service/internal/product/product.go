package product

import (
	"context"

	"github.com/rusalexch/loyalty-group/internal/accrual/internal/common"
)

type productRepo interface {
	Add(ctx context.Context, orderID string, product common.OrderProduct) error
	FindByOrderId(ctx context.Context, orderID string) (common.OrderGoods, error)
}

type productService struct {
	repo productRepo
}

func New(repo productRepo) *productService {
	return &productService{
		repo: repo,
	}
}
