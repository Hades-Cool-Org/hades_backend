package purchase_list

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"hades_backend/app/web/utils/net"
	"net/http"
)

type Router struct {
}

func (u *Router) URL() string {
	return "/lists"
}

const listIdParam = "list_id"
const userIdParam = "userId"

func (u *Router) Router() func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/", u.Create)
		r.Get("/", u.GetAll)
		r.Get("/users/{user_id}", u.GetAllForUser)
		r.Get("/{list_id}", u.Get)
		r.Put("/{list_id}", u.Update)
		r.Delete("/{list_id}", u.Delete)
	}
}

func (u *Router) GetAllForUser(w http.ResponseWriter, r *http.Request) {

	userId := chi.URLParam(r, userIdParam)

	if userId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("userId is empty")))
		return
	}

	//db search
	puchaseList := []*List{
		{
			ID:     "from db",
			UserID: "from db",
			Products: []*Product{
				{
					ID:            "ID_RETORNADO_DO_BANCO",
					Name:          "product1",
					Image:         "url",
					MeasuringUnit: "UN",
				},
				{
					ID:            "ID_RETORNADO_DO_BANCO_2",
					Name:          "product2",
					Image:         "url",
					MeasuringUnit: "UN",
				},
			},
		},
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &GetAllResponse{puchaseList})
}

func (u *Router) GetAll(w http.ResponseWriter, r *http.Request) {

	//db search
	puchaseList := []*List{
		{
			ID:     "from db",
			UserID: "from db",
			Products: []*Product{
				{
					ID:            "ID_RETORNADO_DO_BANCO",
					Name:          "product1",
					Image:         "url",
					MeasuringUnit: "UN",
				},
				{
					ID:            "ID_RETORNADO_DO_BANCO_2",
					Name:          "product2",
					Image:         "url",
					MeasuringUnit: "UN",
				},
			},
		},
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &GetAllResponse{puchaseList})
}

func (u *Router) Get(w http.ResponseWriter, r *http.Request) {
	listId := chi.URLParam(r, listIdParam)

	if listId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("listId is empty")))
		return
	}

	//db search 404 when empty
	list := &List{
		ID:     "from db",
		UserID: "from db",
		Products: []*Product{
			{
				ID:            "ID_RETORNADO_DO_BANCO",
				Name:          "product1",
				Image:         "url",
				MeasuringUnit: "UN",
			},
			{
				ID:            "ID_RETORNADO_DO_BANCO_2",
				Name:          "product2",
				Image:         "url",
				MeasuringUnit: "UN",
			},
		},
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &Response{list})
}

func (u *Router) Create(w http.ResponseWriter, r *http.Request) {

	data := &Request{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	//db search 404 when empty
	list := &List{
		ID:     "from db",
		UserID: "from db",
		Products: []*Product{
			{
				ID:            "ID_RETORNADO_DO_BANCO",
				Name:          "product1",
				Image:         "url",
				MeasuringUnit: "UN",
			},
			{
				ID:            "ID_RETORNADO_DO_BANCO_2",
				Name:          "product2",
				Image:         "url",
				MeasuringUnit: "UN",
			},
		},
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, &Response{list})
}

func (u *Router) Update(w http.ResponseWriter, r *http.Request) {

	data := &Request{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	listId := chi.URLParam(r, listIdParam)

	if listId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("listId is empty")))
		return
	}

	//db update
	//db search 404 when empty
	list := &List{
		ID:     "from db",
		UserID: "from db",
		Products: []*Product{
			{
				ID:            "ID_RETORNADO_DO_BANCO",
				Name:          "product1",
				Image:         "url",
				MeasuringUnit: "UN",
			},
			{
				ID:            "ID_RETORNADO_DO_BANCO_2",
				Name:          "product2",
				Image:         "url",
				MeasuringUnit: "UN",
			},
		},
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &Response{list})
}

func (u *Router) Delete(w http.ResponseWriter, r *http.Request) {

	listIdParam := chi.URLParam(r, listIdParam)

	if listIdParam == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("listIdParam is empty")))
		return
	}

	render.Status(r, http.StatusNoContent)
}
