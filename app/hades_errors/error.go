package hades_errors

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"hades_backend/api/utils/net"
	"hades_backend/app/logging"
	"net/http"
)

var (
	defaultErrorCode = 500

	ErrNotFound = func(entity string) error {
		return NewNotFoundError(errors.New(fmt.Sprintf("%s not found", entity)))
	}

	l = logging.Initialize()
)

func ParseErrResponse(err error) *net.ErrResponse {

	var hadesErr *HadesError

	isHadesErr := errors.As(err, &hadesErr)

	status := func() int {
		if isHadesErr {
			return hadesErr.Status
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
	l.Error("got an error -> "+err.Error(), zap.Int("status", status))
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
