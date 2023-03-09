package net

import (
	"context"
	"github.com/go-chi/render"
	"hades_backend/app/logging"
	"net/http"
)

func RenderError(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	l := logging.FromContext(ctx)
	l.Error(err.Error())
	errResponse := ParseErrResponse(err)
	render.Status(r, errResponse.HTTPStatusCode)
	render.Render(w, r, errResponse)
}
