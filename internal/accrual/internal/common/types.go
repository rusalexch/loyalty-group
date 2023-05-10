package common

// Order начисления на заказ
type Order struct {
	ID      string      `json:"order"`
	Status  OrderStatus `json:"status"`
	Accrual *float64    `json:"accrual,omitempty"`
}

type Reward struct {
	ID     string     `json:"match"`
	Reward float64    `json:"reward"`
	Type   RewardType `json:"reward_type"`
}

type OrderGoods struct {
	ID    string         `json:"order"`
	Goods []OrderProduct `json:"goods"`
}

// OrderProduct
type OrderProduct struct {
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}
