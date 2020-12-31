package utils

import (
	"errors"
	"fmt"
)

func CouldNotParseEmailError(email string) error {
	return errors.New(fmt.Sprintf("unable to parse email %s", email))
}

func InvalidEmailError(email string) error {
	return errors.New(fmt.Sprintf("invalid or malformed email %s", email))
}

func AliasedEmailError(email string) error {
	return errors.New(fmt.Sprintf("invalid email %s, is aliased", email))
}

func ProhibitedEmailError(email string) error {
	return errors.New(fmt.Sprintf("prohibited email %s, domain is disallowed", email))
}

func InsecurePasswordError() error {
	return errors.New("insecure password")
}

func InvalidNameError(name string) error {
	return errors.New(fmt.Sprintf("invalid name %s, contains disallowed characters", name))
}

func SaveUserToDBError(email string) error {
	return errors.New(fmt.Sprintf("error saving user with email %s, reason unknown or user with that email may already exist.", email))
}

func UserNotFoundError(id string) error {
	return errors.New(fmt.Sprintf("user ID %s not found", id))
}
