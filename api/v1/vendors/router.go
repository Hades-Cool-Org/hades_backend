package vendors

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"hades_backend/api/utils/net"
	"hades_backend/app/cmd/vendors"
	vendors2 "hades_backend/app/model"
	"net/http"
	"strconv"
)

type Router struct {
	service *vendors.Service
}

func NewRouter(service *vendors.Service) *Router {
	return &Router{service: service}
}

func (u *Router) URL() string {
	return "/vendors"
}

const vendorIdParam = "vendor_id"

func (u *Router) Router() func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/", u.Create)
		r.Get("/", u.GetAll)
		r.Get("/{vendor_id}", u.Get)
		r.Put("/{vendor_id}", u.Update)
		r.Delete("/{vendor_id}", u.Delete)
	}
}

func (u *Router) GetAll(w http.ResponseWriter, r *http.Request) {

	vs, err := u.service.GetVendors(r.Context())

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &GetAllResponse{vs})
}

func (u *Router) Get(w http.ResponseWriter, r *http.Request) {
	vendorId := chi.URLParam(r, vendorIdParam)

	if vendorId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("vendorId is empty")))
		return
	}

	vendorIdInt, err := strconv.Atoi(vendorId)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("vendorId is not a number: "+err.Error())))
		return
	}

	vendor, err := u.service.GetVendor(r.Context(), uint(vendorIdInt))

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &Response{vendor})
}

func (u *Router) Create(w http.ResponseWriter, r *http.Request) {

	data := &Request{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	vendorId, err := u.service.CreateVendor(r.Context(), data.Vendor)

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	vendor := &vendors2.Vendor{
		ID:       vendorId,
		Name:     data.Name,
		Phone:    data.Phone,
		Type:     data.Type,
		Location: data.Location,
		Contact:  data.Contact,
		Email:    data.Email,
		Cnpj:     data.Cnpj,
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, &Response{vendor})
}

func (u *Router) Update(w http.ResponseWriter, r *http.Request) {

	data := &Request{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	vendorId := chi.URLParam(r, vendorIdParam)

	if vendorId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("vendorId is empty")))
		return
	}

	vendorIdInt, err := strconv.Atoi(vendorId)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("vendorId is not a number: "+err.Error())))
		return
	}

	err = u.service.UpdateVendor(r.Context(), uint(vendorIdInt), data.Vendor)

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	vendor := &vendors2.Vendor{
		ID:       uint(vendorIdInt),
		Name:     data.Name,
		Phone:    data.Phone,
		Type:     data.Type,
		Location: data.Location,
		Contact:  data.Contact,
		Email:    data.Email,
		Cnpj:     data.Cnpj,
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &Response{vendor})
}

func (u *Router) Delete(w http.ResponseWriter, r *http.Request) {

	vendorId := chi.URLParam(r, vendorIdParam)

	if vendorId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("vendorId is empty")))
		return
	}

	vendorIdInt, err := strconv.Atoi(vendorId)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("vendorId is not a number: "+err.Error())))
		return
	}

	err = u.service.DeleteVendor(r.Context(), uint(vendorIdInt))

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
}
