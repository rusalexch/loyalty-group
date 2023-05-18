package gophermart

import "log"

type Config struct {
	Address        string
	DBURL          string
	AccrualAddress string
}

type gophermart struct {
}

func New(conf Config) *gophermart {
	log.Println(conf)
	return &gophermart{}
}

// users	id	int
// 				login	varchar
// 				password	varchar

// user_orders	id	int64
// 							usert_id	int
// 							status	varchar
// 							accrual	float64
// 							uploaded_at	Data

// transaction	id	int
// 							type	varchar
// 							order_id	int64
// 							amount	float64
// 							processed_at	Data
