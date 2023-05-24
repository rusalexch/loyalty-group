package auth

import (
	"context"

	"github.com/rusalexch/loyalty-group/internal/gophermart/internal/app"
)

type userService interface {
	FindByID(ctx context.Context, userID int) (app.User, error)
	FundByLogin(ctx context.Context, login string) (app.User, error)
	Create(ctx context.Context, user app.CreateUser) (app.User, error)
}
