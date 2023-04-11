package net

import (
	"github.com/go-chi/render"
	"net/http"
)

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrForbidden(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 403,
		StatusText:     "Forbidden",
		ErrorText:      err.Error(),
	}
}

type NoContentResponse struct {
}

func (n *NoContentResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NoContent() render.Renderer {
	return &NoContentResponse{}
}
