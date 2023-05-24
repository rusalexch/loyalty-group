package order

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rusalexch/loyalty-group/internal/accrual/internal/app"
)

type order struct {
	ID      string           `db:"id"`
	Status  string          `db:"status"`
	Accrual sql.NullFloat64 `db:"accrual"`
}

type orderRepository struct {
	pool *pgxpool.Pool
}

// New конструктор репозитория заказов
func New(pool *pgxpool.Pool) *orderRepository {
	repo := &orderRepository{
		pool: pool,
	}
	err := repo.init()
	if err != nil {
		log.Println("can't init order repository")
		log.Println(err)
	}

	return repo
}

// init инициализация репозитория заказов
func (repo *orderRepository) init() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := repo.pool.Exec(ctx, sqlCreateOrderStatus)
	if err != nil {
		return err
	}
	_, err = repo.pool.Exec(ctx, sqlCreateOrders)

	return err
}

// Add добавление нового заказа
func (repo *orderRepository) Add(ctx context.Context, orderID string) error {
	_, err := repo.pool.Exec(ctx, sqlAddNewOrder, orderID)

	return err
}

// FindByID поиск заказа по номеру
func (repo *orderRepository) FindByID(ctx context.Context, orderID string) (app.Order, error) {
	var order order
	row := repo.pool.QueryRow(ctx, sqlFindByID, orderID)
	err := row.Scan(&order.ID, &order.Status, &order.Accrual)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return app.Order{}, app.ErrOrderNotFound
	}
	log.Println(err)
	// FIXME cannot scan NULL into *float64

	return dbToJSON(order), err
}

// UpdateStatus изменение статуса заказа
func (repo *orderRepository) UpdateStatus(ctx context.Context, orderID string, status string) error {
	_, err := repo.pool.Exec(ctx, sqlUpdateStatus, orderID, status)

	return err
}

// Update изменения данных заказа
func (repo *orderRepository) Update(ctx context.Context, order app.Order) error {
	_, err := repo.pool.Exec(ctx, sqlUpdate, order.ID, order.Status, order.Accrual)

	return err
}

// Delete удаление заказа
func (repo *orderRepository) Delete(ctx context.Context, orderID string) error {
	_, err := repo.pool.Exec(ctx, sqlDelete, orderID)
	return err
}

// FindRegistered поиск новых заказов для расчета
func (repo *orderRepository) FindRegistered(ctx context.Context) ([]string, error) {
	rows, err := repo.pool.Query(ctx, sqlFindRegistered)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]string, 0)
	for rows.Next() {
		var orderID string
		err := rows.Scan(&orderID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return res, nil
			}
			return nil, err
		}
		res = append(res, orderID)
	}

	return res, nil
}

// dbToJSON преобразование структуры БД в структуру JSON
func dbToJSON(order order) app.Order {
	var accrual *float64 = nil
	if order.Accrual.Valid {
		accrual = &order.Accrual.Float64
	}
	return app.Order{
		ID:      order.ID,
		Status:  order.Status,
		Accrual: accrual,
	}
}
