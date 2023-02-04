package login

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"hades_backend/app/cmd/user"
	"hades_backend/app/web/utils/net"
	"net/http"
)

type Router struct {
	loginService user.Service
}

func NewRouter(loginService *user.Service) *Router {
	return &Router{loginService: *loginService}
}

func (u *Router) URL() string {
	return "/login"
}

func (u *Router) Router() func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/", u.Login)
	}
}

func (u *Router) Login(w http.ResponseWriter, r *http.Request) {
	data := &Request{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}
	login, err := u.loginService.Login(r.Context(), data.Email, data.Password)

	if err != nil {
		render.Render(w, r, net.ErrForbidden(err))
		return
	}
	render.Status(r, http.StatusOK)
	render.Render(w, r, &Response{Token: login.Token})
}
