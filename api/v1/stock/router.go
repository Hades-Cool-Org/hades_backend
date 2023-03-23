package stock

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"hades_backend/api/utils/net"
	stockModel "hades_backend/app/cmd/stock"
	"hades_backend/app/model"
	"net/http"
	"strconv"
	"time"
)

type Router struct {
	service *stockModel.Service
}

func NewRouter(service *stockModel.Service) *Router {
	return &Router{service: service}
}

func (u *Router) URL() string {
	return "/stock"
}

const storeIdParam = "store_id"
const stockIdParam = "stock_id"
const productIdParam = "product_id"

func (u *Router) Router() func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/{stock_id}", u.GetById)
		r.Delete("/{stock_id}", u.Delete)

		r.Post("/store/{store_id}", u.Create)
		r.Get("/store/{store_id}", u.GetByStoreId)
		r.Put("/store/{store_id}/product/{product_id}", u.UpdateProduct)
		r.Post("/store/{store_id}/product/{product_id}", u.AddProducts)
		r.Get("/store/{store_id}/product/{product_id}", u.GetProduct)
		r.Delete("/store/{store_id}/product/{product_id}", u.DeleteProduct)
	}
}

func (u *Router) GetById(w http.ResponseWriter, r *http.Request) {

	stockId := chi.URLParam(r, "stock_id")

	if stockId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("stockId is empty")))
		return
	}

	stockIdInt, err := strconv.Atoi(stockId)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("stockId is not a number: "+err.Error())))
		return
	}

	s, err := u.service.GetStock(r.Context(), uint(stockIdInt))

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &Response{s})
}

func (u *Router) Create(w http.ResponseWriter, r *http.Request) {

	data := &Request{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	storeId := chi.URLParam(r, storeIdParam)

	if storeId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("storeId is empty")))
		return
	}

	createdStockId, err := u.service.CreateStock(r.Context(), data.Stock)

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	//db save
	s := &model.Stock{
		ID:           createdStockId,
		StoreId:      data.StoreId,
		StoreName:    data.StoreName,
		LastModified: time.Now().Format(time.RFC3339Nano),
		Products:     data.Products,
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, &Response{s})
}

func (u *Router) GetByStoreId(w http.ResponseWriter, r *http.Request) {

	storeId := chi.URLParam(r, storeIdParam)

	if storeId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("storeId is empty")))
		return
	}

	storeIdInt, err := strconv.Atoi(storeId)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("storeId is not a number: "+err.Error())))
		return
	}

	s, err := u.service.GetStock(r.Context(), uint(storeIdInt))

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &Response{s})

}

func (u *Router) Delete(w http.ResponseWriter, r *http.Request) {

	stockId := chi.URLParam(r, stockIdParam)

	if stockId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("stockId is empty")))
		return
	}

	stockIdInt, err := strconv.Atoi(stockId)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("stockId is not a number: "+err.Error())))
		return
	}

	err = u.service.DeleteStock(r.Context(), uint(stockIdInt))

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (u *Router) UpdateProduct(w http.ResponseWriter, r *http.Request) {

	data := &ProductRequest{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	storeId := chi.URLParam(r, storeIdParam)

	if storeId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("storeId is empty")))
		return
	}

	storeIdInt, err := strconv.Atoi(storeId)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("storeId is not a number: "+err.Error())))
		return
	}

	productId := chi.URLParam(r, productIdParam)

	if productId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("productId is empty")))
		return
	}

	productIdInt, err := strconv.Atoi(productId)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("productId is not a number: "+err.Error())))
		return
	}

	err = u.service.UpdateProduct(r.Context(), uint(storeIdInt), uint(productIdInt), data.ProductData)

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	//db save
	p := &model.ProductData{
		ProductId:   data.ProductData.ProductId,
		ProductName: data.ProductData.ProductName,
		Current:     data.ProductData.Current,
		Suggested:   data.ProductData.Suggested,
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &ProductResponse{p})

}

func (u *Router) AddProducts(w http.ResponseWriter, r *http.Request) {

	data := &ProductRequestList{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	storeId := chi.URLParam(r, storeIdParam)

	if storeId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("storeId is empty")))
		return
	}

	storeIdInt, err := strconv.Atoi(storeId)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("storeId is not a number: "+err.Error())))
		return
	}

	err = u.service.AddProductToStock(r.Context(), uint(storeIdInt), data.Products)

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
}

func (u *Router) GetProduct(w http.ResponseWriter, r *http.Request) {

	storeId := chi.URLParam(r, storeIdParam)

	if storeId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("storeId is empty")))
		return
	}

	storeIdInt, err := strconv.Atoi(storeId)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("storeId is not a number: "+err.Error())))
		return
	}

	productId := chi.URLParam(r, productIdParam)

	if productId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("productId is empty")))
		return
	}

	productIdInt, err := strconv.Atoi(productId)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("productId is not a number: "+err.Error())))
		return
	}

	p, err := u.service.GetProduct(r.Context(), uint(storeIdInt), uint(productIdInt))

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &ProductResponse{p})
}

func (u *Router) DeleteProduct(w http.ResponseWriter, r *http.Request) {

	storeId := chi.URLParam(r, storeIdParam)

	if storeId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("storeId is empty")))
		return
	}

	storeIdInt, err := strconv.Atoi(storeId)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("storeId is not a number: "+err.Error())))
		return
	}

	productId := chi.URLParam(r, productIdParam)

	if productId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("productId is empty")))
		return
	}

	productIdInt, err := strconv.Atoi(productId)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("productId is not a number: "+err.Error())))
		return
	}

	err = u.service.RemoveProductFromStock(r.Context(), uint(storeIdInt), uint(productIdInt))

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
}
