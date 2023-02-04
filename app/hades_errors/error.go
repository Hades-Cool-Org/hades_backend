package hades_errors

import (
	"errors"
	"fmt"
	"hades_backend/app/web/utils/net"
	"net/http"
)

var (
	defaultErrorCode = 500

	ErrNotFound = func(entity string) error {
		return NewNotFoundError(errors.New(fmt.Sprintf("%s not found", entity)))
	}
)

func ParseErrResponse(err error) *net.ErrResponse {

	status := func() int {
		if e, ok := err.(HadesError); ok {
			return e.Status
		}

		return defaultErrorCode
	}()

	return &net.ErrResponse{
		Err:            err,
		HTTPStatusCode: status,
		StatusText:     err.Error(),
		AppCode:        0,
		ErrorText:      err.Error(),
	}

}

type HadesError struct {
	error
	Status int
}

func NewHadesError(err error, status int) *HadesError {
	if err == nil {
		err = errors.New("unknown error")
	}
	return &HadesError{err, status}
}

func NewNotFoundError(err error) *HadesError {
	return NewHadesError(err, http.StatusNotFound)
}

func NewForbiddenError(err error) *HadesError {
	return NewHadesError(err, http.StatusForbidden)
}
