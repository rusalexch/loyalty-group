package order

import (
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type accrual struct {
	ID      string   `json:"order"`
	Status  string   `json:"status"`
	Accrual *float64 `json:"accrual,omitempty"`
}

type order struct {
	ID         string          `db:"id"`
	UserID     int             `db:"user_id"`
	Status     string          `db:"status"`
	Accrual    sql.NullFloat64 `db:"accrual"`
	UploadedAt time.Time       `db:"uploaded_at"`
}

type updateOrder struct {
	ID      string
	Status  string
	Accrual sql.NullFloat64
}

type Config struct {
	Pool           *pgxpool.Pool
	AccrualAddress string
}

type orderModule struct {
	pool           *pgxpool.Pool
	account        balancer
	tick           *time.Ticker
	accrualAddress string
}
