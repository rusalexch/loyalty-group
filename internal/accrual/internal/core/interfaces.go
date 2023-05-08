package core

import (
	"context"

	"github.com/rusalexch/loyalty-group/internal/accrual/internal/common"
)

type storager interface {
	Add(ctx context.Context, product common.Product) error
	Match(ctx context.Context, match []string) (map[string]common.Product, error)
}