package products

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"hades_backend/web/utils/net"
	"net/http"
)

type Router struct {
}

func (u *Router) URL() string {
	return "/products"
}

const productIdParam = "product_id"

func (u *Router) Router() func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/", u.CreateProduct)
		r.Get("/", u.GetAllProducts)
		r.Get("/{product_id}", u.GetProduct)
		r.Put("/{product_id}", u.UpdateProduct)
		r.Delete("/{product_id}", u.DeleteProduct)
	}
}

func (u *Router) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	//db search
	products := []*Product{
		{
			ID:            "ID_RETORNADO_DO_BANCO",
			Name:          "product1",
			Details:       "details",
			Image:         "url",
			MeasuringUnit: "UN",
		},
		{
			ID:            "ID_RETORNADO_DO_BANCO_2",
			Name:          "product2",
			Details:       "details",
			Image:         "url",
			MeasuringUnit: "UN",
		},
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &AllProductResponse{products})
}

func (u *Router) GetProduct(w http.ResponseWriter, r *http.Request) {
	productId := chi.URLParam(r, productIdParam)

	if productId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("productId is empty")))
		return
	}

	//db search 404 when empty
	product := &Product{
		ID:            productId,
		Name:          "product2",
		Details:       "details",
		Image:         "url",
		MeasuringUnit: "UN",
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &ProductResponse{product})
}

func (u *Router) CreateProduct(w http.ResponseWriter, r *http.Request) {

	data := &ProductRequest{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	//db save
	product := &Product{
		ID:            "ID_RETORNADO_DO_BANCO",
		Name:          data.Name,
		Details:       data.Details,
		Image:         data.Image,
		MeasuringUnit: data.MeasuringUnit,
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, &ProductResponse{product})
}

func (u *Router) UpdateProduct(w http.ResponseWriter, r *http.Request) {

	data := &ProductRequest{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	productId := chi.URLParam(r, productIdParam)

	if productId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("productId is empty")))
		return
	}

	//db update
	product := &Product{
		ID:            productId,
		Name:          data.Name,
		Details:       data.Details,
		Image:         data.Image,
		MeasuringUnit: data.MeasuringUnit,
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &ProductResponse{product})
}

func (u *Router) DeleteProduct(w http.ResponseWriter, r *http.Request) {

	productId := chi.URLParam(r, productIdParam)

	if productId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("productId is empty")))
		return
	}

	render.Status(r, http.StatusOK)
}
