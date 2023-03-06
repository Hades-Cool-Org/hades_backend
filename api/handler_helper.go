package api

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	customMiddleware "hades_backend/api/middleware"
	"hades_backend/api/v1/login"
	"hades_backend/api/v1/product"
	"hades_backend/api/v1/store"
	"hades_backend/api/v1/user"
	"hades_backend/api/v1/vendors"
	productService "hades_backend/app/cmd/product"
	storeService "hades_backend/app/cmd/store"
	userService "hades_backend/app/cmd/user"
	vendorsCmd "hades_backend/app/cmd/vendors"
	productRepository "hades_backend/app/repository/product"
	storeRepository "hades_backend/app/repository/store"
	userRepository "hades_backend/app/repository/user"
	vendorsRepository "hades_backend/app/repository/vendors"
)

type CustomHandler interface {
	Handle(r chi.Router)
}

type MySQLHandler struct {
	DB             *gorm.DB
	userRepository userRepository.Repository
	userService    *userService.Service
	vendorService  *vendorsCmd.Service
	productService *productService.Service
	storeService   *storeService.Service
}

func NewMySQLHandler(db *gorm.DB) *MySQLHandler {
	h := &MySQLHandler{DB: db}

	h.userRepository = userRepository.NewMySqlRepository(db)
	h.userService = userService.NewService(h.userRepository)

	vRepository := vendorsRepository.NewMySqlRepository(db)
	h.vendorService = vendorsCmd.NewService(vRepository)

	pRepository := productRepository.NewMySqlRepository(db)
	h.productService = productService.NewService(pRepository)

	sr := storeRepository.NewMySqlRepository(db)
	h.storeService = storeService.NewService(sr, h.userRepository)

	return h
}

func (m *MySQLHandler) Handle(r chi.Router) {

	loginRouter := m.initLoginRouter()
	r.Route(loginRouter.URL(), loginRouter.Router())

	r.Route("/v1", func(r chi.Router) {

		// Seek, verify and validate JWT tokens
		r.Use(customMiddleware.Verifier(userService.TokenAuth))
		// Handle valid / invalid tokens.
		r.Use(customMiddleware.Authenticator)
		// Extract user
		r.Use(customMiddleware.User)

		userRouter := m.initUserRouter()
		r.Route(userRouter.URL(), userRouter.Router())

		productsRouter := m.initProductRouter()
		r.Route(productsRouter.URL(), productsRouter.Router())

		vendorsRouter := m.initVendorsRouter()
		r.Route(vendorsRouter.URL(), vendorsRouter.Router())

		storeRouter := m.initStoreRouter()
		r.Route(storeRouter.URL(), storeRouter.Router())
	})
}

func (m *MySQLHandler) initLoginRouter() *login.Router {
	return login.NewRouter(m.userService)
}

func (m *MySQLHandler) initUserRouter() *user.Router {
	return user.NewRouter(m.userService)
}

func (m *MySQLHandler) initVendorsRouter() *vendors.Router {
	return vendors.NewRouter(m.vendorService)
}

func (m *MySQLHandler) initProductRouter() *product.Router {
	return product.NewRouter(m.productService)
}

func (m *MySQLHandler) initStoreRouter() *store.Router {
	return store.NewRouter(m.storeService)
}
