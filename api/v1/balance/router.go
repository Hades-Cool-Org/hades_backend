package balance

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"gorm.io/gorm"
	"hades_backend/api/utils/net"
	"hades_backend/app/cmd/balance"
	"net/http"
	"os"
	"strconv"
)

type Router struct {
	db *gorm.DB
}

func NewRouter(db *gorm.DB) *Router {
	err := db.AutoMigrate(&balance.Balance{}, balance.Entry{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return &Router{db: db}
}

func (u *Router) URL() string {
	return "/balance"
}

func (u *Router) Router() func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/{user_id}", u.GetBalance)
	}
}

func (u *Router) GetBalance(w http.ResponseWriter, r *http.Request) {

	userId := chi.URLParam(r, "user_id")

	if userId == "" {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("user_id is empty")))
		return
	}

	userIdInt, err := strconv.Atoi(userId)

	if err != nil {
		render.Render(w, r, net.ErrInvalidRequest(errors.New("user_id is not a number: "+err.Error())))
		return
	}

	bl, err := balance.GetBalance(r.Context(), u.db, uint(userIdInt))

	if err != nil {
		net.RenderError(r.Context(), w, r, err)
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, bl)
}
