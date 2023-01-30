package vendors

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
	//db search
	vendors := []*Vendor{
		{
			ID:       "ID_RETORNADO_DO_BANCO",
			Name:     "Vendor1",
			Email:    "Vendor1@gmail.com",
			Phone:    "+5519999999999",
			Phone2:   "+5519999999999",
			Cnpj:     "23232323232",
			Type:     "CEASA",
			Location: "pavilhão 3, box 18",
		},
		{
			ID:       "ID_RETORNADO_DO_BANCO",
			Name:     "Vendor2",
			Email:    "Vendor2@gmail.com",
			Phone:    "+5519999999999",
			Phone2:   "+5519999999999",
			Cnpj:     "23232323232",
			Type:     "EXTERNO",
			Location: "Rua 30 de julho numero 430",
		},
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &GetAllResponse{vendors})
}

func (u *Router) Get(w http.ResponseWriter, r *http.Request) {
	vendorId := chi.URLParam(r, vendorIdParam)

	if vendorId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("vendorId is empty")))
		return
	}

	//db search 404 when empty
	vendor := &Vendor{
		ID:       vendorId,
		Name:     "Vendor1",
		Email:    "Vendor1@gmail.com",
		Phone:    "+5519999999999",
		Phone2:   "+5519999999999",
		Cnpj:     "23232323232",
		Type:     "CEASA",
		Location: "pavilhão 3, box 18",
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

	//db save
	vendor := &Vendor{
		ID:       "ID_RETORNADO_DO_BANCO",
		Name:     "Vendor1",
		Email:    "Vendor1@gmail.com",
		Phone:    "+5519999999999",
		Phone2:   "+5519999999999",
		Cnpj:     "23232323232",
		Type:     "CEASA",
		Location: "pavilhão 3, box 18",
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

	//db update
	vendor := &Vendor{
		ID:       vendorId,
		Name:     data.Name,
		Email:    data.Email,
		Phone:    data.Phone,
		Phone2:   data.Phone2,
		Cnpj:     data.Cnpj,
		Type:     data.Type,
		Location: data.Location,
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

	render.Status(r, http.StatusNoContent)
}
