package user

import (
	"context"

	"github.com/rusalexch/loyalty-group/internal/gophermart/internal/app"
)

func (um *userModule) FindByID(ctx context.Context, userID int) (app.User, error) {
	return um.findByID(ctx, userID)
}

func (um *userModule) FindByLogin(ctx context.Context, login string) (app.User, error) {
	return um.findByLogin(ctx, login)
}

func (um *userModule) Create(ctx context.Context, user app.CreateUser) (app.User, error) {
	return um.create(ctx, user)
}
