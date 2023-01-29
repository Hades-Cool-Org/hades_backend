package web

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"hades_backend/web/v1/login"
	"hades_backend/web/v1/products"
	"hades_backend/web/v1/users"
	"hades_backend/web/v1/vendors"
	"net/http"
)

func Service() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	userRouter := users.Router{}
	r.Route(userRouter.URL(), userRouter.Router())

	loginRouter := login.Router{}
	r.Route(loginRouter.URL(), loginRouter.Router())

	productsRouter := products.Router{}
	r.Route(productsRouter.URL(), productsRouter.Router())

	vendorsRouter := vendors.Router{}
	r.Route(vendorsRouter.URL(), vendorsRouter.Router())

	return r
}
