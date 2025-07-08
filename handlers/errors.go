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

	// ErrInsecurePassword is returned when the provided password does not meet security requirements.
	ErrInsecurePassword = &HTTPError{
		StatusCode: http.StatusBadRequest,
		Message:    "Password is not secure enough",
	}

	// ErrUserDoesntExist is returned when the user does not exist in database.
	ErrUserDoesntExist = &HTTPError{
		StatusCode: http.StatusNotFound,
		Message:    "User does not exist",
	}

	// ErrIncorrectPassword is returned when the password is incorrect.
	ErrIncorrectPassword = &HTTPError{
		StatusCode: http.StatusUnauthorized,
		Message:    "Incorrect password",
	}

	// ErrInvalidToken is returned when the token is invalid or malformed.
	ErrInvalidToken = &HTTPError{
		StatusCode: http.StatusUnauthorized,
		Message:    "Invalid token",
	}

	// ErrTokenExpired is returned when the token has expired.
	ErrTokenExpired = &HTTPError{
		StatusCode: http.StatusUnauthorized,
		Message:    "Token has expired",
	}
)
