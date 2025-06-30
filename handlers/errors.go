package handlers

import "net/http"

type HTTPError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func (e *HTTPError) Error() string {
	return e.Message
}

var (
	// ErrUserAlreadyExists is returned when a user with the same username already exists.
	ErrUserAlreadyExists = &HTTPError{
		StatusCode: http.StatusConflict,
		Message:    "User already exists",
	}

	// ErrEmailAlreadyExists is returned when a user with the same email already exists.
	ErrEmailAlreadyExists = &HTTPError{
		StatusCode: http.StatusConflict,
		Message:    "Email already exists",
	}
)
