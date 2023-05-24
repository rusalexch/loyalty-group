package auth

import (
	"context"
	"net/http"

	"github.com/rusalexch/loyalty-group/internal/gophermart/internal/app"
)

func (am *authModule) AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			token := r.Header.Get(app.AuthHeader)
			if token == "" {
				next.ServeHTTP(w, r)
				return
			}
			user, err := am.CheckToken(ctx, token)
			if err != nil {
				next.ServeHTTP(w, r)
			}

			r = r.WithContext(context.WithValue(ctx, app.UserKey, user))
			next.ServeHTTP(w, r)
		})
	}
}
