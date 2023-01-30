package stock

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"hades_backend/web/utils/net"
	"net/http"
	"time"
)

type Router struct {
}

func (u *Router) URL() string {
	return "/stock"
}

const storeIdParam = "store_id"
const productIdParam = "product_id"

func (u *Router) Router() func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/{store_id}", u.Create)
		r.Get("/{store_id}", u.Get)
		r.Delete("/{store_id}", u.Delete)
		r.Put("/{store_id}/product/{product_id}", u.UpdateProduct)
		r.Post("/{store_id}/product/{product_id}", u.AddProduct)
		r.Get("/{store_id}/product/{product_id}", u.GetProduct)
		r.Delete("/{store_id}/product/{product_id}", u.DeleteProduct)
	}
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

	//db save
	stock := &Stock{
		StoreId:      storeId,
		LastModified: time.Now().Format(time.RFC3339Nano), //todo
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, &Response{stock})
}

func (u *Router) Get(w http.ResponseWriter, r *http.Request) {

	storeId := chi.URLParam(r, storeIdParam)

	if storeId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("storeId is empty")))
		return
	}

	//db get
	stock := &Stock{
		StoreId:      storeId,
		LastModified: time.Now().Format(time.RFC3339),
		Stock: []*Product{
			{
				ID:        "from db",
				Name:      "from db",
				Current:   5.4,
				Suggested: 2,
			},
			{
				ID:        "from db2",
				Name:      "from db2",
				Current:   5,
				Suggested: 2,
			},
		},
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &Response{stock})
}

func (u *Router) Delete(w http.ResponseWriter, r *http.Request) {

	storeId := chi.URLParam(r, storeIdParam)

	if storeId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("storeId is empty")))
		return
	}

	//db delete
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

	productId := chi.URLParam(r, productIdParam)

	if productId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("productId is empty")))
		return
	}

	//db save
	stock := &Product{
		ID:        productId,
		Name:      "NAME_FROM_DB",
		Current:   data.Current,
		Suggested: data.Suggested,
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &ProductResponse{stock})
}

func (u *Router) AddProduct(w http.ResponseWriter, r *http.Request) {

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

	productId := chi.URLParam(r, productIdParam)

	if productId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("productId is empty")))
		return
	}

	//db save
	stock := &Stock{
		StoreId:      storeId,
		LastModified: time.Now().Format(time.RFC3339),
		Stock: []*Product{ //need to fetch all current?
			{
				ID:        productId,
				Name:      "NAME_FROM_DB",
				Current:   data.Current,
				Suggested: data.Suggested,
			},
		},
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &Response{stock})
}

func (u *Router) GetProduct(w http.ResponseWriter, r *http.Request) {

	storeId := chi.URLParam(r, storeIdParam)

	if storeId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("storeId is empty")))
		return
	}

	productId := chi.URLParam(r, productIdParam)

	if productId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("productId is empty")))
		return
	}

	//db save
	stock := &Product{
		ID:        productId,
		Name:      "NAME_FROM_DB",
		Current:   66.3,
		Suggested: 11.2,
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &ProductResponse{stock})
}

func (u *Router) DeleteProduct(w http.ResponseWriter, r *http.Request) {

	storeId := chi.URLParam(r, storeIdParam)

	if storeId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("storeId is empty")))
		return
	}

	productId := chi.URLParam(r, productIdParam)

	if productId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("productId is empty")))
		return
	}

	//db delete

	render.Status(r, http.StatusNoContent)
}
