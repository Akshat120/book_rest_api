package middleware

import (
	"book_rest_api/internal/config"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func BasicAuthMiddleware(cfg *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return basicAuthHandler(next, cfg)
	}
}

func basicAuthHandler(next http.Handler, cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || username != cfg.BasicAuth.Username || bcrypt.CompareHashAndPassword([]byte(cfg.BasicAuth.Password), []byte(password)) != nil {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
