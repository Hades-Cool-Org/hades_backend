package api

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	customMiddleware "hades_backend/api/middleware"
	"hades_backend/api/v1/login"
	"hades_backend/api/v1/order"
	"hades_backend/api/v1/product"
	"hades_backend/api/v1/purchase_list"
	"hades_backend/api/v1/stock"
	"hades_backend/api/v1/store"
	"hades_backend/api/v1/user"
	"hades_backend/api/v1/vendors"
	productService "hades_backend/app/cmd/product"
	purchaseListService "hades_backend/app/cmd/purchase_list"
	stockService "hades_backend/app/cmd/stock"
	storeService "hades_backend/app/cmd/store"
	userService "hades_backend/app/cmd/user"
	vendorsCmd "hades_backend/app/cmd/vendors"
)

type CustomHandler interface {
	Handle(r chi.Router)
}

type MySQLHandler struct {
	DB                  *gorm.DB
	userRepository      userService.Repository
	userService         *userService.Service
	vendorService       *vendorsCmd.Service
	productService      *productService.Service
	storeService        *storeService.Service
	stockService        *stockService.Service
	purchaseListService *purchaseListService.Service
}

func NewMySQLHandler(db *gorm.DB) *MySQLHandler {
	h := &MySQLHandler{DB: db}

	h.userRepository = userService.NewMySqlRepository(db)
	h.userService = userService.NewService(h.userRepository)

	vRepository := vendorsCmd.NewMySqlRepository(db)
	h.vendorService = vendorsCmd.NewService(vRepository)

	pRepository := productService.NewMySqlRepository(db)
	h.productService = productService.NewService(pRepository)

	sr := storeService.NewMySqlRepository(db)
	h.storeService = storeService.NewService(sr, h.userRepository)

	stockRepository := stockService.NewMySQLRepository(db)
	h.stockService = stockService.NewService(stockRepository)

	purchaseListRepository := purchaseListService.NewRepository(db)
	h.purchaseListService = purchaseListService.NewService(purchaseListRepository)

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

		stockRouter := m.initStockRouter()
		r.Route(stockRouter.URL(), stockRouter.Router())

		purchaseListRouter := m.initPurchaseListRouter()
		r.Route(purchaseListRouter.URL(), purchaseListRouter.Router())

		orderRouter := order.NewRouter(m.DB)
		r.Route(orderRouter.URL(), orderRouter.Router())
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

func (m *MySQLHandler) initStockRouter() *stock.Router {
	return stock.NewRouter(m.stockService)
}

func (m *MySQLHandler) initPurchaseListRouter() *purchase_list.Router {
	return purchase_list.NewRouter(m.purchaseListService)
}
