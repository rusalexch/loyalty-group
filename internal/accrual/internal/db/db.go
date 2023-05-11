package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rusalexch/loyalty-group/internal/accrual/internal/db/internal/order"
	"github.com/rusalexch/loyalty-group/internal/accrual/internal/db/internal/product"
	"github.com/rusalexch/loyalty-group/internal/accrual/internal/db/internal/reward"
)

type database struct {
	OrderRepo   orderRepository
	ProductRepo productRepository
	RewardRepo  rewardRepository
	Store       storager
}

// New конструктор базы данных
func New(dbURL string) *database {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Panic(err)
	}
	orderRepo := order.New(pool)
	productRepo := product.New(pool)
	rewardRepo := reward.New(pool)
	store := NewStore(pool)

	return &database{
		OrderRepo:   orderRepo,
		ProductRepo: productRepo,
		RewardRepo:  rewardRepo,
		Store:       store,
	}
}

type storage struct {
	pool *pgxpool.Pool
}

// конструктор общих методов базы данных
func NewStore(pool *pgxpool.Pool) *storage {
	return &storage{
		pool: pool,
	}
}

// Ping проверка работоспособности БД
func (s *storage) Ping(ctx context.Context) error {
	return s.pool.Ping(ctx)
}
