package order

import "time"

type order struct {
	ID         int64     `json:"number" db:"id"`
	Status     string    `json:"status" db:"status"`
	Accrual    *float64  `json:"accrual,omitempty" db:"accrual"`
	UploadedAt time.Time `json:"uploaded_at" db:"uploaded_at"`
}
