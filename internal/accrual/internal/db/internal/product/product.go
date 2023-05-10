package product

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rusalexch/loyalty-group/internal/accrual/internal/common"
)

type product struct {
	ID          int64   `db:"id"`
	OrderID     string  `db:"order_id"`
	Description string  `db:"description"`
	Price       float64 `db:"price"`
}

type productRepository struct {
	pool *pgxpool.Pool
	mx   *sync.Mutex
}

// New конструктор репозитория товаров
func New(pool *pgxpool.Pool) *productRepository {
	repo := &productRepository{
		pool: pool,
	}

	err := repo.init()
	if err != nil {
		log.Println("can't init product repository")
		log.Fatal(err)
	}

	return repo
}

// init создание таблицы если еще не создана
func (repo *productRepository) init() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := repo.pool.Exec(ctx, sqlCreateProducts)

	return err
}

func (repo *productRepository) Add(ctx context.Context, orderID string, product common.OrderProduct) error {
	repo.mx.Lock()
	defer repo.mx.Unlock()

	_, err := repo.pool.Exec(ctx, sqlAddProduct, orderID, product.Description, product.Price)

	return err
}

func (repo *productRepository) FindByOrderId(ctx context.Context, orderID string) (order common.OrderGoods, err error) {
	repo.mx.Lock()
	defer repo.mx.Unlock()

	rows, err := repo.pool.Query(ctx, sqlFindOrderProducts, orderID)
	if err != nil {
		return common.OrderGoods{}, err
	}
	defer rows.Close()

	goods := make([]product, 0, 10)

	for rows.Next() {
		var product product
		err = rows.Scan(&product.ID, &product.OrderID, &product.Description, &product.Price)
		if err != nil {
			return
		}
		goods = append(goods, product)
	}

	err = rows.Err()
	if err != nil {
		return
	}

	order.ID = orderID
	order.Goods = make([]common.OrderProduct, 0, len(goods))
	for _, g := range goods {
		order.Goods = append(order.Goods, dbToJSON(g))
	}
	return
}

func dbToJSON(product product) common.OrderProduct {
	return common.OrderProduct{
		Description: product.Description,
		Price:       product.Price,
	}
}
