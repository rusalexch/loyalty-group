package order

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"
	"github.com/rusalexch/loyalty-group/internal/accrual/internal/common"
)

type order struct {
	ID      string             `db:"id"`
	Status  common.OrderStatus `db:"status"`
	Accrual float64            `db:"accrual"`
}

type orderRepository struct {
	pool *pgxpool.Pool
	mx   *sync.Mutex
}

func New(pool *pgxpool.Pool) *orderRepository {
	repo := &orderRepository{
		pool: pool,
	}
	err := repo.init()
	if err != nil {
		log.Println("can't init order repository")
		log.Panic(err)
	}

	return repo
}

func (repo *orderRepository) init() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = repo.pool.Exec(ctx, sqlCreateOrderStatus)
	if err != nil {
		return err
	}
	_, err = repo.pool.Exec(ctx, sqlCreateOrders)

	return
}

func (repo *orderRepository) Add(ctx context.Context, orderID string) error {
	repo.mx.Lock()
	defer repo.mx.Unlock()

	_, err := repo.pool.Exec(ctx, sqlAddNewOrder, orderID)

	return err
}

func (repo *orderRepository) FindById(ctx context.Context, orderID string) (common.Order, error) {
	repo.mx.Lock()
	defer repo.mx.Unlock()

	var order order
	row := repo.pool.QueryRow(ctx, sqlFindByID)
	err := row.Scan(order)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return common.Order{}, common.ErrOrderNotFound
	}

	return dbToJSON(order), err
}

func (repo *orderRepository) UpdateStatus(ctx context.Context, orderID string, status common.OrderStatus) error {
	repo.mx.Lock()
	defer repo.mx.Unlock()

	_, err := repo.pool.Exec(ctx, sqlUpdateStatus, orderID, status)

	return err
}

func (repo *orderRepository) Update(ctx context.Context, order common.Order) error {
	repo.mx.Lock()
	defer repo.mx.Unlock()

	_, err := repo.pool.Exec(ctx, sqlUpdate, order.ID, order.Status, order.Accrual)

	return err
}

func dbToJSON(order order) common.Order {
	return common.Order{
		ID:      order.ID,
		Status:  order.Status,
		Accrual: &order.Accrual,
	}
}
