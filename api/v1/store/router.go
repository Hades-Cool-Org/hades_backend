package store

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"hades_backend/api/utils/net"
	storeService "hades_backend/app/cmd/store"
	storeModel "hades_backend/app/model/store"
	"net/http"
	"strconv"
)

type Router struct {
	service *storeService.Service
}

func NewRouter(service *storeService.Service) *Router {
	return &Router{service: service}
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
		r.Delete("/{store_id}/couriers", u.RemoveCourier)
		r.Post("/{store_id}/couriers", u.AddCourier)
	}
}

func (u *Router) GetAll(w http.ResponseWriter, r *http.Request) {
	stores, err := u.service.GetAllStores(r.Context())

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
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

	storeIdInt, err := strconv.Atoi(storeId)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("storeId is not a number: "+err.Error())))
		return
	}

	s, err := u.service.GetStore(r.Context(), uint(storeIdInt))

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

	//db insert
	storeId, err := u.service.CreateStore(r.Context(), data.Store)

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	s := &storeModel.Store{
		ID:       storeId,
		Name:     data.Name,
		Address:  data.Address,
		User:     data.User,
		Couriers: data.Couriers,
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, &Response{s})
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

	storeIdInt, err := strconv.Atoi(storeId)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("storeId is not a number: "+err.Error())))
		return
	}

	data.Store.ID = uint(storeIdInt)

	err = u.service.UpdateStore(r.Context(), data.Store)

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	//db update
	s := &storeModel.Store{
		ID:       uint(storeIdInt),
		Name:     data.Name,
		Address:  data.Address,
		User:     data.User,
		Couriers: data.Couriers,
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &Response{s})
}

func (u *Router) Delete(w http.ResponseWriter, r *http.Request) {

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

	err = u.service.DeleteStore(r.Context(), uint(storeIdInt))

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (u *Router) RemoveCourier(w http.ResponseWriter, r *http.Request) {

	data := &UpdateCouriersRequest{}

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

	err = u.service.RemoveCouriers(r.Context(), uint(storeIdInt), data.Couriers)

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (u *Router) AddCourier(w http.ResponseWriter, r *http.Request) {

	data := &UpdateCouriersRequest{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	//TODO: we should definitively avoid duplicating this code. Abstraction with render doesnt work well in go
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

	err = u.service.AddCouriers(r.Context(), uint(storeIdInt), data.Couriers)
	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	//db save
	render.Status(r, http.StatusNoContent)
}
