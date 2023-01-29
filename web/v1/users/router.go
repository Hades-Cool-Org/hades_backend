package users

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
	return "/users"
}

const userIdParam = "user_id"

func (u *Router) Router() func(r chi.Router) {
	return func(r chi.Router) {
		r.Route("/", func(r chi.Router) {
			//TODO CHECK ROLE
			//r.Use(ArticleCtx)
			//Load obj  on the request context https://github.com/go-chi/chi/blob/16a24da68ae7311e8191d92c67597e5530c5817e/_examples/rest/main.go#L323
			r.Post("/", u.CreateUser)
			r.Put("/{user_id}", u.UpdateUser)
			r.Delete("/{user_id}", u.DeleteUser)
		})
	}
}

func (u *Router) CreateUser(w http.ResponseWriter, r *http.Request) {

	data := &UserRequest{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	//db save
	user := &User{
		ID:    "ID_RETORNADO_DO DB",
		Name:  data.Name,
		Email: data.Email,
		Phone: data.Phone,
		Roles: data.Roles,
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, &UserResponse{user})
}

func (u *Router) UpdateUser(w http.ResponseWriter, r *http.Request) {

	data := &UserRequest{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	userId := chi.URLParam(r, userIdParam)

	if userId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("userId is empty")))
		return
	}

	//db update user
	user := &User{
		ID:    userId,
		Name:  data.Name,
		Email: data.Email,
		Phone: data.Phone,
		Roles: data.Roles,
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &UserResponse{user})
}

func (u *Router) DeleteUser(w http.ResponseWriter, r *http.Request) {

	userId := chi.URLParam(r, userIdParam)

	if userId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("userId is empty")))
		return
	}

	//db delete user

	render.Status(r, http.StatusOK)
}
