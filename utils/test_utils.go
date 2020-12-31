package utils

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func AssertErrorsEqual(t *testing.T, expectedErr error, actualError error) {
	if expectedErr != nil && !cmp.Equal(expectedErr, actualError, cmpopts.EquateErrors()) {
		if actualError == nil {
			t.Errorf("Got no error when '%s' was expected", expectedErr)
		} else if expectedErr.Error() != actualError.Error() {
			t.Errorf("Expected err to be '%s' but got '%s'", expectedErr, actualError)
		}
	} else if expectedErr == nil && actualError != nil {
		t.Errorf("Expected no error but got '%v'", actualError)
	}
}
