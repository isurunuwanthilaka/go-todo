package errors

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type RestError struct {
	Message      string `json:"message"`
	Status       int    `json:"status"`
	ErrorMessage error  `json:"-"`
}

func (r *RestError) Error() string {
	return fmt.Sprintf("Error message -> %v", r.ErrorMessage)
}

func NewError(message string, status int, error string) *RestError {
	return &RestError{
		Message:      message,
		Status:       status,
		ErrorMessage: errors.New(error),
	}
}

func NewBadRequestError(message string) *RestError {
	return &RestError{
		Message:      message,
		Status:       http.StatusBadRequest,
		ErrorMessage: errors.New("bad_request"),
	}
}
