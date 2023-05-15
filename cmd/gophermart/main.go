package main

import (
	"github.com/rusalexch/loyalty-group/internal/config"
	"github.com/rusalexch/loyalty-group/internal/gophermart"
)

func main() {
	conf := config.NewGophermartConfig()

	gophermart.New(conf)
}
