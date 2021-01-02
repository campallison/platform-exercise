package utils

import (
	"errors"
	"fmt"
	"net/http"
)

type APIError struct {
	Message string
	Errors  interface{}
	Code    int
}

func (a APIError) Error() string {
	return a.Message
}

func NewAPIError(message string, errors interface{}, code int) error {
	return APIError{message, errors, code}
}

func CouldNotParseEmailError(email string) error {
	return NewAPIError(
		fmt.Sprintf("unable to parse email %s", email),
		errors.New("could not parse email"),
		http.StatusUnprocessableEntity,
	)
}

func InvalidEmailError(email string) error {
	return NewAPIError(
		fmt.Sprintf("invalid or malformed email %s", email),
		errors.New("invalid email provided"),
		http.StatusBadRequest,
	)
}

func AliasedEmailError(email string) error {
	return NewAPIError(
		fmt.Sprintf("invalid email %s, is aliased", email),
		errors.New(""),
		http.StatusBadRequest,
	)
}

func ProhibitedEmailError(email string) error {
	return NewAPIError(
		fmt.Sprintf("prohibited email %s, domain is disallowed", email),
		errors.New("email domain prohibited"),
		http.StatusBadRequest,
	)
}

func InsecurePasswordError() error {
	return NewAPIError(
		"password does not meet strength threshold",
		errors.New("insecure password provided"),
		http.StatusBadRequest,
	)
}

func InvalidNameError(name string) error {
	return NewAPIError(
		fmt.Sprintf("invalid name %s, contains disallowed characters", name),
		errors.New("name contains invalid characters"),
		http.StatusBadRequest,
	)
}

func SaveUserToDBError(email string) error {
	return NewAPIError(
		fmt.Sprintf("error saving user with email %s, reason unknown or user with that email may already exist.", email),
		errors.New("error saving user to database"),
		http.StatusBadRequest,
	)
}

func UserNotFoundError(id string) error {
	return NewAPIError(
		fmt.Sprintf("user ID %s not found", id),
		errors.New("user not found by ID"),
		http.StatusBadRequest,
	)
}

func LoginFailedError() error {
	return APIError{
		Message: "login failed",
		Errors:  errors.New("login failed"),
		Code:    http.StatusInternalServerError,
	}
}

func InvalidTokenError() error {
	return APIError{
		Message: "invalid token",
		Errors:  errors.New("invalid token"),
		Code:    http.StatusUnauthorized,
	}
}

func TokenCheckFailedError() error {
	return APIError{
		Message: "invalid token check failed",
		Errors:  errors.New("invalid token check failed, please retry"),
		Code:    http.StatusInternalServerError,
	}
}

func TokenSignatureError() error {
	return APIError{
		Message: "token not signed as expected",
		Errors:  errors.New("token not signed as expected"),
		Code:    http.StatusUnauthorized,
	}
}

func ParseTokenError(err error) error {
	return APIError{
		Message: "unable to parse token",
		Errors:  err,
		Code:    http.StatusUnauthorized,
	}
}

func AuthHeaderError() error {
	return APIError{
		Message: "missing authorization header",
		Errors:  errors.New("missing authorization header"),
		Code:    http.StatusUnauthorized,
	}
}
