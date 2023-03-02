package middleware

import (
	"github.com/go-chi/jwtauth/v5"
	"hades_backend/app/environment"
	"net/http"
)

func Verifier(ja *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	if environment.IsProd() {
		return jwtauth.Verify(ja, jwtauth.TokenFromHeader, jwtauth.TokenFromCookie)
	}

	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(hfn)
	}
}

func Authenticator(next http.Handler) http.Handler {
	if environment.IsProd() {
		return jwtauth.Authenticator(next)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
