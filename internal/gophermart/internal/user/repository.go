package user

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rusalexch/loyalty-group/internal/gophermart/internal/app"
)

func (um *userModule) initRepo() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := um.pool.Exec(ctx, sqlCreateTable)
	if err != nil {
		log.Println("user > init > can't create table")
		log.Fatal(err)
	}

}

func (um *userModule) findByID(ctx context.Context, userID int) (app.User, error) {
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

func (um *userModule) findByLogin(ctx context.Context, login string) (app.User, error) {
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

func (um *userModule) create(ctx context.Context, user app.CreateUser) (app.User, error) {
	_, err := um.pool.Exec(ctx, sqlAdd, user.Login, user.Password)
	if err != nil {
		return app.User{}, err
	}

	return um.findByLogin(ctx, user.Login)
}
