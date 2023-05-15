package product

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rusalexch/loyalty-group/internal/accrual/internal/app"
)

const insertStmt = "insertProduct"

type product struct {
	ID          int64   `db:"id"`
	OrderID     int64   `db:"order_id"`
	Description string  `db:"description"`
	Price       float64 `db:"price"`
}

type productRepository struct {
	pool *pgxpool.Pool
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

// init инициализация репозитория товаров
func (repo *productRepository) init() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := repo.pool.Exec(ctx, sqlCreateProducts)

	return err
}

// Add добавление новых товаров
func (repo *productRepository) Add(ctx context.Context, orderID int64, product []app.OrderProduct) error {
	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(ctx, insertStmt, sqlAddProduct)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	for _, p := range product {
		_, err := repo.pool.Exec(ctx, stmt.Name, orderID, p.Description, p.Price)
		if err != nil {
			tx.Rollback(ctx)
		}
	}

	tx.Commit(ctx)
	return nil
}

// FindByOrderID поиск товаров по номеру заказа
func (repo *productRepository) FindByOrderID(ctx context.Context, orderID int64) (order app.OrderGoods, err error) {
	rows, err := repo.pool.Query(ctx, sqlFindOrderProducts, orderID)
	if err != nil {
		return app.OrderGoods{}, err
	}
	defer rows.Close()

	order.ID = orderID
	order.Goods = make([]app.OrderProduct, 0)

	for rows.Next() {
		var product product
		err = rows.Scan(&product.ID, &product.OrderID, &product.Description, &product.Price)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return order, nil
			}
			return
		}
		order.Goods = append(order.Goods, dbToJSON(product))
	}

	return
}

// dbToJSON преобразование структуры БД в структуру JSON
func dbToJSON(product product) app.OrderProduct {
	return app.OrderProduct{
		Description: product.Description,
		Price:       product.Price,
	}
}
