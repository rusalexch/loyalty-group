package order

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rusalexch/loyalty-group/internal/gophermart/app"
)

func (om *orderModule) initRepository() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := om.pool.Exec(ctx, sqlCreateTable)
	if err != nil {
		log.Println("order > init > can't create user_orders table")
		log.Fatal(err)
	}
}

func (om *orderModule) findByID(ctx context.Context, orderID string) (order, error) {
	var ord order
	row := om.pool.QueryRow(ctx, sqlFindByID, orderID)
	err := row.Scan(&ord.ID, &ord.UserID, &ord.Status, &ord.Accrual, &ord.UploadedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return order{}, app.ErrNotFound
		}
		return order{}, err
	}

	return ord, nil
}

func (om *orderModule) findByUserID(ctx context.Context, userID int) ([]order, error) {
	orders := make([]order, 0, 10)
	rows, err := om.pool.Query(ctx, sqlFundByUserID, userID)
	if err != nil {
		return []order{}, err
	}

	for rows.Next() {
		var ord order
		err := rows.Scan(&ord.ID, &ord.UserID, &ord.Status, &ord.Accrual, &ord.UploadedAt)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []order{}, app.ErrNotFound
			}
			return []order{}, err
		}
		orders = append(orders, ord)
	}

	return orders, nil
}

func (om *orderModule) add(ctx context.Context, userID int, orderID string) error {
	_, err := om.pool.Exec(ctx, sqlAdd, orderID, userID)

	return err
}

func (om *orderModule) findRegistered(ctx context.Context) ([]order, error) {
	orders := make([]order, 0, 10)
	rows, err := om.pool.Query(ctx, sqlFindRegistered)
	if err != nil {
		return []order{}, err
	}
	for rows.Next() {
		var ord order
		err := rows.Scan(&ord.ID, &ord.UserID, &ord.Status, &ord.Accrual, &ord.UploadedAt)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []order{}, app.ErrNotFound
			}
			return []order{}, err
		}
		orders = append(orders, ord)
	}

	return orders, nil
}

func (om *orderModule) updateOrder(ctx context.Context, updOrder updateOrder) error {
	_, err := om.pool.Exec(ctx, sqlUpdateOrder, updOrder.ID, updOrder.Status, updOrder.Accrual)

	return err
}

func dbToJSON(ord order) app.Order {
	var accrual *float64 = nil
	if ord.Accrual.Valid {
		accrual = &ord.Accrual.Float64
	}
	return app.Order{
		ID:         ord.ID,
		Status:     ord.Status,
		Accrual:    accrual,
		UploadedAt: ord.UploadedAt,
	}
}
