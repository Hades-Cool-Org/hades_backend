package user

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"hades_backend/app/cmd/user"
	"hades_backend/app/hades_errors"
	"hades_backend/app/web/utils/net"
	"net/http"
)

type Router struct {
	UserService *user.Service
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
			r.Post("/", u.Create)
			r.Put("/{user_id}", u.Update)
			r.Delete("/{user_id}", u.Delete)
		})
	}
}

func (u *Router) Create(w http.ResponseWriter, r *http.Request) {

	data := &Request{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	//db save
	err := u.UserService.CreateUser(r.Context(), data.ToModel())

	if err != nil {
		errResponse := hades_errors.ParseErrResponse(err)
		render.Status(r, errResponse.HTTPStatusCode)
		render.Render(w, r, errResponse)
		return
	}

	user := &User{
		ID:         "ID_RETORNADO_DO DB",
		Name:       data.Name,
		Email:      data.Email,
		Phone:      data.Phone,
		Roles:      data.Roles,
		FirstLogin: true,
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, &Response{user})
}

func (u *Router) Update(w http.ResponseWriter, r *http.Request) {

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

	//db update user
	user := &User{
		ID:    userId,
		Name:  data.Name,
		Email: data.Email,
		Phone: data.Phone,
		Roles: data.Roles,
	}

	err := u.UserService.UpdateUser(r.Context(), userId, user.ToModel())

	if err != nil {
		errResponse := hades_errors.ParseErrResponse(err)
		render.Status(r, errResponse.HTTPStatusCode)
		render.Render(w, r, errResponse)
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &Response{user})
}

func (u *Router) Delete(w http.ResponseWriter, r *http.Request) {

	userId := chi.URLParam(r, userIdParam)

	if userId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("userId is empty")))
		return
	}

	//db delete user

	render.Status(r, http.StatusNoContent)
}