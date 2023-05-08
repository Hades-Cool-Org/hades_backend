package occurrence

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"hades_backend/api/utils/net"
	"net/http"
	"time"
)

type Router struct {
}

func (u *Router) URL() string {
	return "/occurrences"
}

const occurrenceId = "occurrence_id"

func (u *Router) Router() func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/", u.GetAll)
		r.Delete("/{occurrence_id}", u.Delete)
		r.Get("/{occurrence_id}", u.Get)
	}
}

func (u *Router) Get(w http.ResponseWriter, r *http.Request) {
	boxId := chi.URLParam(r, occurrenceId)

	if boxId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("boxId is empty")))
		return
	}

	delivery := &Occurrence{
		ID:        "idfromdb",
		State:     "created",
		StartDate: time.Now().Format(time.RFC3339),
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &Response{delivery})
}

func (u *Router) Delete(w http.ResponseWriter, r *http.Request) {
	occurrenceId := chi.URLParam(r, occurrenceId)

	if occurrenceId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("occurrenceId is empty")))
		return
	}

	//db delete

	render.Status(r, http.StatusNoContent)
}

func (u *Router) GetAll(w http.ResponseWriter, r *http.Request) {

	userId := chi.URLParam(r, userIdParam)

	delivery := &Occurrence{
		ID:        "idfromdb",
		State:     "created",
		User:      &User{ID: userId},
		StartDate: time.Now().Format(time.RFC3339),
	}

	deliveryList := []*Occurrence{
		delivery,
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &ListResponse{deliveryList})
}
