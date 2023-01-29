package users

import (
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

func (u *Router) Router() func(r chi.Router) {
	return func(r chi.Router) {
		r.Route("/", func(r chi.Router) {
			//TODO CHECK ROLE
			//r.Use(ArticleCtx)            // Load the *Article on the request context
			r.Post("/", u.CreateUser) // POST /articles
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
