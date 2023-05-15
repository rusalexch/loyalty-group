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
