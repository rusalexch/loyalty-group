package account

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rusalexch/loyalty-group/internal/gophermart/app"
)

type transaction struct {
	ID          int       `json:"-" db:"id"`
	Type        string    `json:"-" db:"transaction_type"`
	OrderID     int64     `json:"order" db:"order_id"`
	Amount      float64   `json:"sum" db:"amount"`
	ProcessedAt time.Time `json:"processed_at" db:"processed_at"`
}

func (am *accountModule) createTable() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := am.pool.Exec(ctx, sqlCreateTable)
	if err != nil {
		log.Println("account > create table > can't create table")
		log.Fatalln(err)
	}
}

func (am *accountModule) addDebit(ctx context.Context, orderID int64, amount float64) error {
	_, err := am.pool.Exec(ctx, sqlAddDebit, orderID, amount)

	return err
}

func (am *accountModule) addCredit(ctx context.Context, orderID int64, amount float64) error {
	_, err := am.pool.Exec(ctx, sqlAddCredit, orderID, amount)

	return err
}

func (am *accountModule) currentBalance(ctx context.Context, userID int) (balance, error) {
	var debit, credit float64

	row := am.pool.QueryRow(ctx, sqlGetUserCurrentBalance, userID)
	err := row.Scan(debit, credit)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return balance{}, app.ErrNotFound
		}
		return balance{}, err
	}

	return balance{
		Current:  debit - credit,
		Withdraw: credit,
	}, nil
}

func (am *accountModule) userCredit(ctx context.Context, userId int) ([]transaction, error) {
	rows, err := am.pool.Query(ctx, sqlGetUserCredit, userId)
	if err != nil {
		return []transaction{}, err
	}
	credits := make([]transaction, 0, 10)

	for rows.Next() {
		var tr transaction
		err = rows.Scan(&tr.ID, &tr.Type, &tr.OrderID, &tr.Amount, &tr.ProcessedAt)
		if err != nil {
			return []transaction{}, err
		}

		credits = append(credits, tr)
	}

	return credits, nil
}
