package config

import (
	"flag"
	"log"
	"os"

	"github.com/rusalexch/loyalty-group/internal/accrual"
	"github.com/rusalexch/loyalty-group/internal/gophermart"
)

const (
	defaultAddr           = "127.0.0.1:8080"
	defaultDBURL          = ""
	defaultAccrualAddress = "127.0.0.1:8081"
)

var (
	address        *string
	dbURL          *string
	accrualAddress *string
)

func init() {
	address = flag.String("a", defaultAddr, "set address")
	dbURL = flag.String("d", defaultDBURL, "set address")
	accrualAddress = flag.String("r", defaultAccrualAddress, "set address")
}

func NewAccrualConfig() accrual.Config {
	flag.Parse()
	parseENV()

	if *dbURL == "" || *address == "" {
		log.Panic("accrual > config > db url and address is required")
	}

	return accrual.Config{
		Address: *address,
		DBURL:   *dbURL,
	}
}

func NewGophermartConfig() gophermart.Config {
	flag.Parse()
	parseENV()

	if *dbURL == "" || *address == "" || *accrualAddress == "" {
		log.Panic("gophermart > config > db url and address is required")
	}

	return gophermart.Config{
		Address:        *address,
		DBURL:          *dbURL,
		AccrualAddress: *accrualAddress,
	}
}

func parseENV() {
	if val, ok := os.LookupEnv("RUN_ADDRESS"); ok {
		address = &val
	}
	if val, ok := os.LookupEnv("DATABASE_URI"); ok {
		dbURL = &val
	}
	if val, ok := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS"); ok {
		accrualAddress = &val
	}
}
