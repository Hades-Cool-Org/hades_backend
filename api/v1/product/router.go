package product

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"hades_backend/api/utils/net"
	"hades_backend/app/cmd/product"
	productModel "hades_backend/app/model/product"
	"net/http"
	"strconv"
)

type Router struct {
	service *product.Service
}

func NewRouter(service *product.Service) *Router {
	return &Router{service: service}
}

func (u *Router) URL() string {
	return "/products"
}

const productIdParam = "product_id"

func (u *Router) Router() func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/", u.Create)
		r.Get("/", u.GetAll)
		r.Get("/{product_id}", u.Get)
		r.Put("/{product_id}", u.Update)
		r.Delete("/{product_id}", u.Delete)
	}
}

func (u *Router) GetAll(w http.ResponseWriter, r *http.Request) {
	//db search
	products, err := u.service.GetProducts(r.Context())

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &GetAllResponse{products})
}

func (u *Router) Get(w http.ResponseWriter, r *http.Request) {
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

	p, err := u.service.GetProduct(r.Context(), uint(productIdInt))

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &Response{p})
}

func (u *Router) Create(w http.ResponseWriter, r *http.Request) {

	data := &Request{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	//db save
	p, err := u.service.CreateProduct(r.Context(), data.Product)

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	productResp := &productModel.Product{
		ID:            p,
		Name:          data.Name,
		Details:       data.Details,
		Image:         data.Image,
		MeasuringUnit: data.MeasuringUnit,
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, &Response{productResp})
}

func (u *Router) Update(w http.ResponseWriter, r *http.Request) {

	data := &Request{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
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

	id := uint(productIdInt)
	err = u.service.UpdateProduct(r.Context(), id, data.Product)

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	productResp := &productModel.Product{
		ID:            id,
		Name:          data.Name,
		Details:       data.Details,
		Image:         data.Image,
		MeasuringUnit: data.MeasuringUnit,
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &Response{productResp})
}

func (u *Router) Delete(w http.ResponseWriter, r *http.Request) {

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

	err = u.service.DeleteProduct(r.Context(), uint(productIdInt))

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
}
