package common

type Reward struct {
	ID     string     `json:"match"`
	Reward float64    `json:"reward"`
	Type   RewardType `json:"reward_type"`
}

type OrderProduct struct {
	Name  string  `json:"description"`
	Price float64 `json:"price"`
}

type Order struct {
	ID      string      `json:"order"`
	Status  OrderStatus `json:"status"`
	Accrual *float64    `json:"accrual,omitempty"`
}
