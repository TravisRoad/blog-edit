package middleware

import (
	"net/http"
	"strings"

	"github.com/TravisRoad/blog-edit/global"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if len(authHeader) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if strings.TrimPrefix(authHeader, "Bearer ") != global.Config.AuthToken {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
