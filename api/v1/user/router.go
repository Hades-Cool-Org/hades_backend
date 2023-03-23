package user

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"hades_backend/api/middleware"
	"hades_backend/api/utils/net"
	"hades_backend/app/cmd/user"
	userModel "hades_backend/app/model"
	"net/http"
	"strconv"
)

type Router struct {
	UserService *user.Service
}

func NewRouter(service *user.Service) *Router {
	return &Router{UserService: service}
}

func (r2 *Router) URL() string {
	return "/users"
}

const userIdParam = "user_id"

func (r2 *Router) Router() func(r chi.Router) {
	return func(r chi.Router) {
		r.Route("/", func(r chi.Router) {

			r.Group(func(r chi.Router) {
				r.Use(middleware.PermissionCheck(middleware.FnIsAdmin))
				r.Post("/", r2.Create)
			})

			r.Group(func(r chi.Router) {
				r.Use(middleware.PermissionCheck(middleware.FnCheckUserIdOrAdmin))
				r.Put("/{user_id}", r2.Update)
				r.Delete("/{user_id}", r2.Delete)
			})
		})
	}
}

func (r2 *Router) Create(w http.ResponseWriter, r *http.Request) {

	data := &Request{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	//db save
	id, err := r2.UserService.CreateUser(r.Context(), data.User)

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	u := &userModel.User{
		ID:         id,
		Name:       data.Name,
		Email:      data.Email,
		Phone:      data.Phone,
		Roles:      data.Roles,
		FirstLogin: true,
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, &Response{u})
}

func (r2 *Router) Update(w http.ResponseWriter, r *http.Request) {

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

	userIdInt, err := strconv.Atoi(userId)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("userId is not a number: "+err.Error())))
		return
	}

	//db update u
	u := &userModel.User{
		ID:    uint(userIdInt),
		Name:  data.Name,
		Email: data.Email,
		Phone: data.Phone,
		Roles: data.Roles,
	}

	err = r2.UserService.UpdateUser(r.Context(), uint(userIdInt), u)

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &Response{u})
}

func (r2 *Router) Delete(w http.ResponseWriter, r *http.Request) {

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

	//db delete user
	err = r2.UserService.DeleteUser(r.Context(), uint(userIdInt))

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (r2 *Router) Get(w http.ResponseWriter, r *http.Request) {

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

	//db get user
	u, err := r2.UserService.GetUser(r.Context(), uint(userIdInt))

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &Response{u})
}
