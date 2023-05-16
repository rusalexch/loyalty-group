package accrual

import (
	"log"

	"github.com/rusalexch/loyalty-group/internal/accrual/internal/db"
	"github.com/rusalexch/loyalty-group/internal/accrual/internal/handlers"
	"github.com/rusalexch/loyalty-group/internal/accrual/internal/service"
)

type Config struct {
	Address string
	DBURL   string
}

func Start(config Config) {
	storage := db.New(config.DBURL)
	srv := service.New(service.ServiceConfig{
		Store:       storage.Store,
		OrderRepo:   storage.OrderRepo,
		ProductRepo: storage.ProductRepo,
		RewardRepo:  storage.RewardRepo,
	})
	server := handlers.New(config.Address, srv)

	log.Printf("accrual > starting on %s\n", config.Address)

	server.Start()
}
