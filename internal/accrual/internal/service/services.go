package service

import (
	"context"

	"github.com/rusalexch/loyalty-group/internal/accrual/internal/common"
)

type ServiceConfig struct {
	Store       storager
	OrderRepo   orderRepository
	ProductRepo productRepository
	RewardRepo  rewardRepository
}

type service struct {
	store       storager
	orderRepo   orderRepository
	productRepo productRepository
	rewardRepo  rewardRepository
}

func New(conf ServiceConfig) *service {
	return &service{
		store:       conf.Store,
		orderRepo:   conf.OrderRepo,
		productRepo: conf.ProductRepo,
		rewardRepo:  conf.RewardRepo,
	}
}

func (s *service) Ping(ctx context.Context) error {
	return s.store.Ping(ctx)
}

func (s *service) GetOrder(ctx context.Context, orderID string) (float64, error) {
	return 0, nil
}
func (s *service) AddProduct(ctx context.Context, product common.Reward) error {
	return nil
}

func (s *service) AddOrder(ctx context.Context, order common.OrderGoods) error {
	return nil
}

// func (c *core) Add(ctx context.Context, product common.Reward) error {
// 	return c.store.Add(ctx, product)
// }

// func (c *core) Calc(ctx context.Context, goods []common.OrderProduct) (float64, error) {
// 	matchers := make([]string, 0, len(goods))
// 	for _, g := range goods {
// 		matchers = append(matchers, g.Description)
// 	}
// 	productMap, err := c.store.Match(ctx, matchers)
// 	if err != nil {
// 		return 0, err
// 	}
// 	var accrual float64
// 	for _, g := range goods {
// 		if acc, ok := productMap[g.Description]; ok {
// 			accrual += calc(g.Price, acc)
// 		}
// 	}
// 	return accrual, nil
// }

// func calc(price float64, acc common.Reward) float64 {
// 	switch acc.Type {
// 	case common.Percentage:
// 		return price * acc.Reward / 100
// 	case common.Fixed:
// 		return acc.Reward
// 	default:
// 		return 0
// 	}
// }
