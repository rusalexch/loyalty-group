package app

import "github.com/rusalexch/loyalty-group/internal/validator"

// Order начисления на заказ
type Order struct {
	ID      string   `json:"order"`
	Status  string   `json:"status"`
	Accrual *float64 `json:"accrual,omitempty"`
}

type Reward struct {
	ID     string  `json:"match"`
	Reward float64 `json:"reward"`
	Type   string  `json:"reward_type"`
}

type OrderGoods struct {
	ID    string         `json:"order"`
	Goods []OrderProduct `json:"goods"`
}

func (o *OrderGoods) IsValid() bool {
	return validator.IsValid(o.ID)
}

// OrderProduct
type OrderProduct struct {
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}
