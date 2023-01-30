package store

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
	return "/store"
}

const storeIdParam = "store_id"

func (u *Router) Router() func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/", u.Create)
		r.Get("/", u.GetAll)
		r.Get("/{store_id}", u.Get)
		r.Put("/{store_id}", u.Update)
		r.Delete("/{store_id}", u.Delete)
	}
}

func (u *Router) GetAll(w http.ResponseWriter, r *http.Request) {
	//db search
	stores := []*Store{
		{
			ID:      "ID_RETORNADO_DO_BANCO",
			Name:    "store1",
			Address: "Rua 30 de julho, 545",
		},
		{
			ID:      "ID_RETORNADO_DO_BANCO",
			Name:    "store2",
			Address: "Rua 30 de julho, 545",
		},
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &GetAllResponse{stores})
}

func (u *Router) Get(w http.ResponseWriter, r *http.Request) {
	storeId := chi.URLParam(r, storeIdParam)

	if storeId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("storeId is empty")))
		return
	}

	//db search 404 when empty
	store := &Store{
		ID:      storeId,
		Name:    "store2",
		Address: "Address",
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &Response{store})
}

func (u *Router) Create(w http.ResponseWriter, r *http.Request) {

	data := &Request{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	//db save
	store := &Store{
		ID:      "ID_RETORNADO_DO_BANCO",
		Name:    data.Name,
		Address: data.Address,
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, &Response{store})
}

func (u *Router) Update(w http.ResponseWriter, r *http.Request) {

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

	//db update
	store := &Store{
		ID:      storeId,
		Name:    data.Name,
		Address: data.Address,
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &Response{store})
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
