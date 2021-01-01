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
		errors.New("user not found"),
		http.StatusBadRequest,
	)
}
