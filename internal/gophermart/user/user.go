package user

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rusalexch/loyalty-group/internal/gophermart/app"
)

type Config struct {
	Pool *pgxpool.Pool
}

type userModule struct {
	pool *pgxpool.Pool
}

func New(conf Config) *userModule {
	module := &userModule{
		pool: conf.Pool,
	}
	module.init()

	return module
}

func (um *userModule) init() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := um.pool.Exec(ctx, sqlCreateTable)
	if err != nil {
		log.Println("user > init > can't create table")
		log.Fatal(err)
	}
}

func (um *userModule) FindByID(ctx context.Context, userID int) (app.User, error) {
	var user app.User
	row := um.pool.QueryRow(ctx, sqlFindByID, userID)
	err := row.Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return app.User{}, app.ErrNotFound
		}
		return app.User{}, err
	}

	return user, nil
}

func (um *userModule) FundByLogin(ctx context.Context, login string) (app.User, error) {
	var user app.User
	row := um.pool.QueryRow(ctx, sqlFindByLogin, login)
	err := row.Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return app.User{}, app.ErrNotFound
		}
		return app.User{}, err
	}

	return user, nil
}

func (um *userModule) Create(ctx context.Context, user app.CreateUser) (app.User, error) {
	_, err := um.pool.Exec(ctx, sqlAdd, user.Login, user.Password)
	if err != nil {
		return app.User{}, err
	}

	return um.FundByLogin(ctx, user.Login)
}
