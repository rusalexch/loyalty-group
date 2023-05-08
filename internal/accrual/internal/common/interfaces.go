package common

import "context"

type Accrualer interface {
	Calc(ctx context.Context, products []OrderProduct) (uint, error)
}

type GoodsCreater interface {
	Add(ctx context.Context, product Reward) error
}
