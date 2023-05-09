package core

import (
	"context"

	"github.com/rusalexch/loyalty-group/internal/accrual/internal/common"
)

type core struct {
	store storager
}

func New(store storager) *core {
	return &core{
		store: store,
	}
}

func (c *core) Ping(ctx context.Context) error {
	return c.store.Ping(ctx)
}

func (c *core) Add(ctx context.Context, product common.Reward) error {
	return c.store.Add(ctx, product)
}

func (c *core) Calc(ctx context.Context, goods []common.OrderProduct) (float64, error) {
	matchers := make([]string, 0, len(goods))
	for _, g := range goods {
		matchers = append(matchers, g.Name)
	}
	productMap, err := c.store.Match(ctx, matchers)
	if err != nil {
		return 0, err
	}
	var accrual float64
	for _, g := range goods {
		if acc, ok := productMap[g.Name]; ok {
			accrual += calc(g.Price, acc)
		}
	}
	return accrual, nil
}

func calc(price float64, acc common.Reward) float64 {
	switch acc.Type {
	case common.Percentage:
		return price * acc.Reward / 100
	case common.Fixed:
		return acc.Reward
	default:
		return 0
	}
}
