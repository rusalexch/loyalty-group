package main

import (
	"log"

	"github.com/rusalexch/loyalty-group/internal/accrual"
)

func main() {
	conf := accrual.Config{}

	err := accrual.Start(conf)
	if err != nil {
		log.Panic(err)
	}
}
