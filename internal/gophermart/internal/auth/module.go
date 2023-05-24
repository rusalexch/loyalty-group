package auth

import (
	"net/http"

	"github.com/rusalexch/loyalty-group/internal/gophermart/internal/app"
)

func New(conf Config) *authModule {
	module := &authModule{
		userService: conf.UserService,
		jwtSecret:   conf.JwtSecret,
	}

	return module
}

func (am *authModule) Middlewares() []app.Middleware {
	return []app.Middleware{
		am.createAuthMiddleware(),
	}
}

func (am *authModule) Handlers() []app.Handler {
	register := app.Handler{
		Method:  http.MethodPost,
		Pattern: "/api/user/register",
		Handler: am.register,
	}
	login := app.Handler{
		Method:  http.MethodPost,
		Pattern: "/api/user/login",
		Handler: am.login,
	}

	return []app.Handler{register, login}
}
