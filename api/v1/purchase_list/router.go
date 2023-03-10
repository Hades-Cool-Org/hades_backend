package purchase_list

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"hades_backend/api/utils/net"
	purchaseListModel "hades_backend/app/cmd/purchase_list"
	"hades_backend/app/model/purchase_list"
	"net/http"
	"strconv"
)

type Router struct {
	service *purchaseListModel.Service
}

func NewRouter(service *purchaseListModel.Service) *Router {
	return &Router{service: service}
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

	userIdInt, err := strconv.Atoi(userId)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("userId is not a number: "+err.Error())))
		return
	}

	lists, err := u.service.GetPurchaseListsByUserId(r.Context(), uint(userIdInt))
	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &GetAllResponse{lists})
}

func (u *Router) GetAll(w http.ResponseWriter, r *http.Request) {

	lists, err := u.service.GetPurchaseLists(r.Context())
	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &GetAllResponse{lists})
}

func (u *Router) Get(w http.ResponseWriter, r *http.Request) {

	listId := chi.URLParam(r, listIdParam)

	if listId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("listId is empty")))
		return
	}

	listIdInt, err := strconv.Atoi(listId)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("listId is not a number: "+err.Error())))
		return
	}

	list, err := u.service.GetPurchaseList(r.Context(), uint(listIdInt))
	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
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

	id, err := u.service.CreatePurchaseList(r.Context(), data.PurchaseList)
	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	list := &purchase_list.PurchaseList{ID: id}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, &Response{list})
}

func (u *Router) Update(w http.ResponseWriter, r *http.Request) {

	listIdInPAram := chi.URLParam(r, listIdParam)

	if listIdInPAram == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("listIdInPAram is empty")))
		return
	}

	listId, err := strconv.Atoi(listIdInPAram)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("listIdInPAram is not a number: "+err.Error())))
		return
	}

	data := &Request{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	data.PurchaseList.ID = uint(listId)

	err = u.service.UpdatePurchaseList(r.Context(), uint(listId), data.PurchaseList)
	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (u *Router) Delete(w http.ResponseWriter, r *http.Request) {

	listIdInPAram := chi.URLParam(r, listIdParam)

	if listIdInPAram == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("listIdInPAram is empty")))
		return
	}

	listId, err := strconv.Atoi(listIdInPAram)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("listIdInPAram is not a number: "+err.Error())))
		return
	}

	err = u.service.DeletePurchaseList(r.Context(), uint(listId))
	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
}
