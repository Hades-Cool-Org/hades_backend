package web

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"go.uber.org/zap"
	"hades_backend/app/cmd/auth"
	customMiddleware "hades_backend/app/web/middleware"
	"hades_backend/app/web/v1/login"
	"hades_backend/app/web/v1/product"
	"hades_backend/app/web/v1/user"
	"hades_backend/app/web/v1/vendors"
	"net/http"
)

func Service(l *zap.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(customMiddleware.Logger(l))
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	loginRouter := login.Router{}
	r.Route(loginRouter.URL(), loginRouter.Router())

	r.Group(func(r chi.Router) {

		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(auth.TokenAuth))

		// Handle valid / invalid tokens. In this example, we use
		// the provided authenticator middleware, but you can write your
		// own very easily, look at the Authenticator method in jwtauth.go
		// and tweak it, its not scary.
		r.Use(jwtauth.Authenticator)

		r.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			w.Write([]byte(fmt.Sprintf("protected area. hi %v", claims)))
		})

		r.Route("/v1", func(r chi.Router) {

			userRouter := user.Router{}
			r.Route(userRouter.URL(), userRouter.Router())

			productsRouter := product.Router{}
			r.Route(productsRouter.URL(), productsRouter.Router())

			vendorsRouter := vendors.Router{}
			r.Route(vendorsRouter.URL(), vendorsRouter.Router())
		})
	})

	return r
}
