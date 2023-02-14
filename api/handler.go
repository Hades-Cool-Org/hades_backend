package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"go.uber.org/zap"
	customMiddleware "hades_backend/api/middleware"
	"hades_backend/api/utils"
	"hades_backend/api/v1/login"
	"hades_backend/api/v1/product"
	"hades_backend/api/v1/user"
	"hades_backend/api/v1/vendors"
	user2 "hades_backend/app/cmd/user"
	vendorsCmd "hades_backend/app/cmd/vendors"
	"hades_backend/app/database"
	user3 "hades_backend/app/repository/user"
	vendorsRepository "hades_backend/app/repository/vendors"
	"net/http"
)

var (
	db             = database.DB
	userRepository = user3.NewMySqlRepository(db)
	userService    = user2.NewService(userRepository)
	vendorService  = vendorsCmd.NewService(vendorsRepository.NewMySqlRepository(db))
)

func Handler(l *zap.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(customMiddleware.Logger(l))
	r.Use(middleware.Recoverer)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	loginRouter := initLoginRouter()
	r.Route(loginRouter.URL(), loginRouter.Router())

	r.Group(func(r chi.Router) {

		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(user2.TokenAuth))
		// Handle valid / invalid tokens.
		r.Use(jwtauth.Authenticator)
		// Extract user
		r.Use(customMiddleware.User)

		r.Route("/v1", func(r chi.Router) {

			userRouter := initUserRouter()
			r.Route(userRouter.URL(), userRouter.Router())

			r.Route(productsRouter.URL(), productsRouter.Router())

			vendorsRouter := initVendorsRouter()
			r.Route(vendorsRouter.URL(), vendorsRouter.Router())
		})
	})

	utils.GenerateDocs(r)
	return r
}

func initLoginRouter() *login.Router {
	return login.NewRouter(userService)
}

func initUserRouter() *user.Router {
	return user.NewRouter(userService)
}

func initVendorsRouter() *vendors.Router {
	return vendors.NewRouter(vendorService)
}
