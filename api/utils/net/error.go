package net

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"hades_backend/app/logging"
	"net/http"
)

var (
	defaultErrorCode = 500
)

func ParseErrResponse(err error) *Error {

	var hadesErr *HadesError

	isHadesErr := errors.As(err, &hadesErr)

	status := func() int {
		if isHadesErr {
			return hadesErr.Status
		}

		return defaultErrorCode
	}()

	return &Error{
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

func NewHadesError(ctx context.Context, err error, status int) *HadesError {
	l := logging.FromContext(ctx)
	l.Error("got an error -> "+err.Error(), zap.Int("status", status))
	if err == nil {
		err = errors.New("unknown error")
	}
	return &HadesError{err, status}
}

func NewForbiddenError(ctx context.Context, err error) *HadesError {
	return NewHadesError(ctx, err, http.StatusForbidden)
}
