package main

import (
	"github.com/rusalexch/loyalty-group/internal/accrual"
)

func main() {
	conf := accrual.Config{}

	accrual.Start(conf)

}
