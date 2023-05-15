package main

import (
	"github.com/rusalexch/loyalty-group/internal/accrual"
	"github.com/rusalexch/loyalty-group/internal/config"
)

func main() {
	conf := config.NewAccrualConfig()

	accrual.Start(conf)

}
