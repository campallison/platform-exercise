package utils

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func AssertErrorsEqual(t *testing.T, err1 error, err2 error) {
	if err1 != nil && !cmp.Equal(err1, err2, cmpopts.EquateErrors()) {
		if err2 == nil {
			t.Errorf("Got error '%s' when no error was expected", err1)
		} else if err1.Error() != err2.Error() {
			t.Errorf("Expected err to be '%s' but was '%s'", err1, err2)
		}
	} else if err1 == nil && err2 != nil {
		t.Errorf("Got no error when error was expected to be '%v'", err2)
	}
}
