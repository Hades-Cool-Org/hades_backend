package middleware

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"hades_backend/api/utils/net"
	"hades_backend/app/identity"
	"net/http"
	"strconv"
)

var (
	extractUserIDFromRequest = func(r *http.Request) uint {
		userId, _ := strconv.Atoi(chi.URLParam(r, "user_id"))
		return uint(userId)
	}

	FnCheckUserIdOrAdmin = func(r *http.Request) bool {
		id := identity.FromContext(r.Context())
		return id.IsAdmin() || id.UserId == extractUserIDFromRequest(r)
	}

	FnIsAdmin = func(r *http.Request) bool {
		id := identity.FromContext(r.Context())
		return id.IsAdmin()
	}
)

func PermissionCheck(fn func(r *http.Request) bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			if !fn(r) {
				e := net.NewForbiddenError(r.Context(), errors.New("you are not allowed perform this operation"))
				errResponse := net.ParseErrResponse(e)
				render.Status(r, errResponse.HTTPStatusCode)
				render.Render(w, r, errResponse)
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
