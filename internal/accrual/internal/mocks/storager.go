package mocks

import (
	"context"
	"sync"

	"github.com/rusalexch/loyalty-group/internal/accrual/internal/common"
)

type MockStorager struct {
	goods map[string]common.Reward
	mx    sync.Mutex
}

func MockAccrualStoragerNew() *MockStorager {
	return &MockStorager{
		goods: make(map[string]common.Reward),
	}
}

func (mas *MockStorager) Ping(ctx context.Context) error {
	return nil
}

func (mas *MockStorager) Add(ctx context.Context, product common.Reward) error {
	mas.mx.Lock()
	defer mas.mx.Unlock()

	mas.goods[product.ID] = product

	return nil
}

func (mas *MockStorager) Match(ctx context.Context, match []string) (map[string]common.Reward, error) {
	mas.mx.Lock()
	defer mas.mx.Unlock()

	res := make(map[string]common.Reward)
	for _, key := range match {
		if v, ok := mas.goods[key]; ok {
			res[key] = v
		}
	}

	return res, nil
}
