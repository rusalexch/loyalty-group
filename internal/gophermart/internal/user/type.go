package user

import "github.com/jackc/pgx/v5/pgxpool"

type Config struct {
	Pool *pgxpool.Pool
}

type userModule struct {
	pool *pgxpool.Pool
}
