package accrual

import (
	"log"

	"github.com/rusalexch/loyalty-group/internal/accrual/internal/core"
	"github.com/rusalexch/loyalty-group/internal/accrual/internal/db"
)

type Config struct {
	Address string
	DBURL   string
}

func Start(config Config) {
	stor := db.New(config.DBURL)
	c := core.New(stor)
	log.Println(c)
}
