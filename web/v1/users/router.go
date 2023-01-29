package users

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"hades_backend/cmd/users/auth"
	"hades_backend/web/utils/net"
	"net/http"
)

type UserRouter struct {
	LoginService auth.Service
}

func (u *UserRouter) URL() string {
	return "/users"
}

func (u *UserRouter) Router() func(r chi.Router) {
	return func(r chi.Router) {

		r.Route("/admin", func(r chi.Router) {
			//TODO CHECK ROLE
			//r.Use(ArticleCtx)            // Load the *Article on the request context
			r.Post("/", u.CreateUser) // POST /articles
		})

		r.Post("/login", u.Login)
	}
}
func (u *UserRouter) Login(w http.ResponseWriter, r *http.Request) {

	data := &UserLoginRequest{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, net.ErrInvalidRequest(err))
		return
	}

	login, err := u.LoginService.Login(data.Email, data.Password)

	if err != nil {
		render.Render(w, r, net.ErrForbidden(err))

	} else {
		render.Status(r, http.StatusOK)
		render.Render(w, r, &UserLoginResponse{Token: login.Token})
	}
}

func (u *UserRouter) CreateUser(w http.ResponseWriter, r *http.Request) {

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
