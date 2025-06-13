package middleware

import (
	"book_rest_api/internal/config"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func BasicAuthMiddleware(next http.Handler) http.Handler {
	config := config.GetConfig()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || username != config.BasicAuth.Username || bcrypt.CompareHashAndPassword([]byte(config.BasicAuth.Password), []byte(password)) != nil {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
