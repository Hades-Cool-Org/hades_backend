package middleware

import (
	"fmt"
	"github.com/go-chi/jwtauth/v5"
	"hades_backend/app/identity"
	"net/http"
)

// User get's user email from JWT token and adds it to the request context.
func User(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())

		userId, _ := claims["user_id"].(float64)

		var roles []string

		for _, item := range claims["roles"].([]interface{}) {
			roles = append(roles, fmt.Sprintf("%v", item))
		}

		r = r.WithContext(identity.WithUser(r.Context(), uint(userId), roles))

		next.ServeHTTP(w, r)
	})
}
