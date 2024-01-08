package middleware

import (
	"log/slog"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		slog.Info(r.RequestURI, slog.String("method", r.Method))
	})
}
