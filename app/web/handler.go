package web

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
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

	r.Route("/v1", func(r chi.Router) {
		userRouter := user.Router{}
		r.Route(userRouter.URL(), userRouter.Router())

		productsRouter := product.Router{}
		r.Route(productsRouter.URL(), productsRouter.Router())

		vendorsRouter := vendors.Router{}
		r.Route(vendorsRouter.URL(), vendorsRouter.Router())
	})

	loginRouter := login.Router{}
	r.Route(loginRouter.URL(), loginRouter.Router())

	return r
}
