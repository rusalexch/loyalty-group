package accrual

import (
	"github.com/rusalexch/loyalty-group/internal/accrual/internal/core"
	"github.com/rusalexch/loyalty-group/internal/accrual/internal/db"
	"github.com/rusalexch/loyalty-group/internal/accrual/internal/handlers"
)

type Config struct {
	Address string
	DBURL   string
}

func Start(config Config) {
	storage := db.New(config.DBURL)
	service := core.New(storage)
	server := handlers.New(config.Address, service)

	server.Start()
}
