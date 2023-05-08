package mocks

import (
	"context"
	"sync"

	"github.com/rusalexch/loyalty-group/internal/accrual/internal/common"
)

type MockAccrualStorager struct {
	goods map[string]common.Product
	mx    sync.Mutex
}

func MockAccrualStoragerNew() *MockAccrualStorager {
	return &MockAccrualStorager{
		goods: make(map[string]common.Product),
	}
}

func (mas *MockAccrualStorager) Add(ctx context.Context, product common.Product) error {
	mas.mx.Lock()
	defer mas.mx.Unlock()

	mas.goods[product.ID] = product

	return nil
}

func (mas *MockAccrualStorager) Match(ctx context.Context, match []string) (map[string]common.Product, error) {
	mas.mx.Lock()
	defer mas.mx.Unlock()

	res := make(map[string]common.Product)
	for _, key := range match {
		if v, ok := mas.goods[key]; ok {
			res[key] = v
		}
	}

	return res, nil
}
