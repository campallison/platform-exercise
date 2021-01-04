package utils

import (
	"os"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
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

func CreateTestToken(userID string, email string) string {
	expiry := time.Now().In(time.UTC).Add(time.Hour * 12)
	unsignedToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS512"), jwt.MapClaims{
		"Id":        userID,
		"ExpiresAt": expiry,
		"Subject":   email,
	})

	signedToken, _ := unsignedToken.SignedString([]byte(os.Getenv("SigningSecret")))
	return signedToken
}

func CreateTestAuthHeader(token string, contentType string) map[string]string {
	tokenValue := "bearer " + token

	return map[string]string{
		"Content-Type":  contentType,
		"Authorization": tokenValue,
	}
}
