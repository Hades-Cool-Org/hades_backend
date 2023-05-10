package conference

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"hades_backend/api/utils/net"
	"hades_backend/app/cmd/conference"
	"hades_backend/app/cmd/occurence"
	"hades_backend/app/cmd/user"
	"hades_backend/app/model"
	"net/http"
	"strconv"
	"time"
)

type Router struct {
}

func (u *Router) URL() string {
	return "/conference"
}

const occurrenceIdParam = "occurrence_id"

func (u *Router) Router() func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/", u.Conference)

		r.Get("/occurrences", u.GetAll)
		r.Delete("/occurrences/{occurrence_id}", u.Delete)
		r.Get("/occurrences/{occurrence_id}", u.GetByOccurrenceId)
	}
}

func (u *Router) Conference(w http.ResponseWriter, r *http.Request) {

	var request Request

	if err := render.Bind(r, &request); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	err := conference.DoConference(r.Context(), request.Occurrence)
	if err != nil {
		net.RenderError(r.Context(), w, r, err)
	}

	render.Status(r, http.StatusOK)
}

func (u *Router) GetByOccurrenceId(w http.ResponseWriter, r *http.Request) {
	occurrenceId := chi.URLParam(r, occurrenceIdParam)

	if occurrenceId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("occurrenceId is empty")))
		return
	}

	occurrenceIdInt, err := strconv.Atoi(occurrenceId)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("occurrenceId is not a number: "+err.Error())))
		return
	}

	occurrence, err := occurence.GetOccurrence(r.Context(), uint(occurrenceIdInt))

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("occurrenceId is not a number: "+err.Error())))
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &Response{entityToModel(occurrence)})
}

func (u *Router) Delete(w http.ResponseWriter, r *http.Request) {
	occurrenceId := chi.URLParam(r, occurrenceIdParam)

	if occurrenceId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("occurrenceId is empty")))
		return
	}

	occurrenceIdInt, err := strconv.Atoi(occurrenceId)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("occurrenceId is not a number: "+err.Error())))
		return
	}

	err = occurence.DeleteOccurrence(r.Context(), uint(occurrenceIdInt))

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
	}

	render.Status(r, http.StatusNoContent)
}

func (u *Router) GetAll(w http.ResponseWriter, r *http.Request) {

	opts := &occurence.GetOccurrenceOptions{
		Params: r.URL.Query(),
	}

	occurrences, err := occurence.GetOccurrences(r.Context(), opts)

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &ListResponse{Occurrences: entitiesToModels(occurrences)})
}

func entitiesToModels(occurrences []*occurence.Occurrence) []*model.Occurrence {
	var models []*model.Occurrence
	for _, occurrence := range occurrences {
		models = append(models, entityToModel(occurrence))
	}
	return models
}

func entityToModel(oc *occurence.Occurrence) *model.Occurrence {

	fnUser := func(user *user.User) *model.User {
		if user == nil {
			return nil
		}
		return &model.User{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Phone: user.Phone,
		}
	}

	fnItems := func(items []*occurence.Item) []*model.OccurrenceItem {
		var models []*model.OccurrenceItem
		for _, item := range items {
			models = append(models, &model.OccurrenceItem{
				ProductID:     item.ProductID,
				Type:          model.ToOccurrenceType(item.Type),
				Name:          item.Name,
				MeasuringUnit: item.MeasuringUnit,
				Quantity:      item.Quantity,
			})
		}
		return models
	}

	deletedAtStr := ""

	if oc.DeletedAt.Valid {
		deletedAtStr = oc.DeletedAt.Time.Format(time.RFC3339)
	}

	return &model.Occurrence{
		ID:          oc.ID,
		DeliveryID:  oc.DeliveryID,
		StoreID:     oc.StoreID,
		User:        fnUser(oc.User),
		Items:       fnItems(oc.Items),
		CreatedDate: oc.CreatedAt.Format(time.RFC3339),
		EndDate:     deletedAtStr,
	}
}
