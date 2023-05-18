package account

import "time"

type transaction struct {
	ID          int       `json:"-" db:"id"`
	Type        string    `json:"-" db:"transaction_type"`
	OrderID     int64     `json:"order" db:"order_id"`
	Amount      float64   `json:"sum" db:"amount"`
	ProcessedAt time.Time `json:"processed_at" db:"processed_at"`
}

type balance struct {
	Current  float64 `json:"current"`
	Withdraw float64 `json:"withdrawn"`
}

type withdrawRequest struct {
	OrderID int64   `json:"order"`
	Amount  float64 `json:"sum"`
}
