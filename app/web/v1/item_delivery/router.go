package item_delivery

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"hades_backend/app/web/utils/net"
	"net/http"
	"time"
)

type Router struct {
}

func (u *Router) URL() string {
	return "/item"
}

const userIdParam = "user_id"
const itemIdParam = "item_id"

func (u *Router) Router() func(r chi.Router) {
	return func(r chi.Router) {

		r.Post("/", u.Create) //end user turn
		r.Get("/", u.GetAll)
		r.Delete("/{item_id}", u.Delete)
		r.Get("/{item_id}", u.Get)
		r.Post("/{item_id}", u.Complete)
	}
}

func (u *Router) Create(w http.ResponseWriter, r *http.Request) {

	data := &Request{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	userId := chi.URLParam(r, userIdParam)

	if userId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("userId is empty")))
		return
	}

	data.User.ID = userId

	delivery := &Delivery{
		ID:        "idfromdb",
		State:     "created",
		User:      data.User,
		StartDate: time.Now().Format(time.RFC3339),
		EndDate:   "",
		Items:     data.Items,
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &Response{delivery})
}

func (u *Router) Complete(w http.ResponseWriter, r *http.Request) {
	itemId := chi.URLParam(r, itemIdParam)

	if itemId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("item is empty")))
		return
	}

	render.Status(r, http.StatusOK)
}

func (u *Router) Get(w http.ResponseWriter, r *http.Request) {
	boxId := chi.URLParam(r, itemIdParam)

	if boxId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("boxId is empty")))
		return
	}

	delivery := &Delivery{
		ID:        boxId,
		State:     "created",
		StartDate: time.Now().Format(time.RFC3339),
		EndDate:   "",
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &Response{delivery})
}

func (u *Router) Delete(w http.ResponseWriter, r *http.Request) {
	boxId := chi.URLParam(r, itemIdParam)

	if boxId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("boxId is empty")))
		return
	}

	//db delete

	render.Status(r, http.StatusNoContent)
}

func (u *Router) GetAll(w http.ResponseWriter, r *http.Request) {

	userId := chi.URLParam(r, userIdParam)

	delivery := &Delivery{
		ID:    "boxId",
		State: "created",
		User: &User{
			ID:   userId,
			Name: "",
		},
		StartDate: time.Now().Format(time.RFC3339),
		EndDate:   "",
	}

	deliveryList := []*Delivery{
		delivery,
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &ListResponse{deliveryList})
}
