package api

import (
	"encoding/json"
	"errors"
	"net/http"
)

type CommonError struct {
	Code int
	Err  error
}

func (s CommonError) Error() string {
	return s.Err.Error()
}

func (s CommonError) Status() int {
	return s.Code
}

func (s CommonError) Response() []byte {
	var (
		e map[string]interface{} = map[string]interface{}{
			"message": s.Error(),
			"status":  s.Code,
		}
		r []byte
	)

	r, _ = json.Marshal(e)
	return r
}

func (s CommonError) WriteResponse(w http.ResponseWriter) {
	w.WriteHeader(s.Status())
	w.Write(s.Response())
}

func NewBadRequestErr(err error) *CommonError {
	return &CommonError{
		Code: http.StatusBadRequest,
		Err:  errors.New(http.StatusText(http.StatusBadRequest)),
	}
}

func NewNotFoundErr(err error) *CommonError {
	return &CommonError{
		Code: http.StatusNotFound,
		Err:  errors.New(http.StatusText(http.StatusNotFound)),
	}
}

func NewServiceUnavailableErr(err error) *CommonError {
	return &CommonError{
		Code: http.StatusServiceUnavailable,
		Err:  errors.New(http.StatusText(http.StatusServiceUnavailable)),
	}
}

func NotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err = NewNotFoundErr(nil)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		err.WriteResponse(w)
	}
}
