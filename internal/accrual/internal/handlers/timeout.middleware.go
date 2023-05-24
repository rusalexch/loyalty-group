package handlers

import (
	"context"
	"net/http"
	"time"
)

func timeoutMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer func() {
			cancel()
			if ctx.Err() == context.DeadlineExceeded {
				w.WriteHeader(http.StatusGatewayTimeout)
			}
		}()
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
