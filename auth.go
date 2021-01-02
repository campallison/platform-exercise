package platform_exercise

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/campallison/platform-exercise/utils"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Credential struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password", validate:"required"`
}

func (c Credential) CheckPassword(hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(c.Password))
	return err == nil
}

func Login(creds Credential) (LoginResponse, error) {
	db := Init()
	var user User
	var response LoginResponse

	if err := db.Table("users").Where("email = ?", creds.Email).First(&user).Error; err != nil {
		return response, utils.LoginFailedError()
	}

	if creds.CheckPassword(user.Password) {
		expiry := time.Now().In(time.UTC).Add(time.Hour * 12)
		unsignedToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS512"), jwt.MapClaims{
			"Id":        user.ID,
			"ExpiresAt": expiry,
			"Subject":   user.Email,
		})

		signedToken, err := unsignedToken.SignedString([]byte(os.Getenv("signing_secret")))
		if err != nil {
			return response, utils.LoginFailedError()
		}

		response.AccessToken = signedToken
		response.Expiry = expiry

		return response, nil
	}

	return LoginResponse{}, utils.LoginFailedError()
}

func Logout(authHeader string) (LogoutResponse, error) {
	tokenString, err := getTokenFromAuthHeader(authHeader)
	if err != nil {
		return LogoutResponse{}, err
	}

	token := InvalidToken{
		Token: tokenString,
	}

	db := Init()
	deleteStaleInvalidTokens()
	if err := db.Save(&token); err != nil {
		return LogoutResponse{}, utils.APIError{
			Message: "logout failed",
			Errors:  err,
			Code:    http.StatusInternalServerError,
		}
	}

	return LogoutResponse{Success: true}, nil
}

func deleteStaleInvalidTokens() {
	db := Init()
	if err := db.Table(
		"invalid_tokens",
	).Delete(
		&InvalidToken{},
	).Where(
		"created_at < now() - interval '12 hours'",
	).Error; err != nil {
		return
	}
}

func CheckToken(authHeader string, userID string) error {
	tokenString, err := getTokenFromAuthHeader(authHeader)
	if err != nil {
		return err
	}

	invalidTokenCheck := isInvalidToken(tokenString)
	if invalidTokenCheck == nil {
		return utils.TokenCheckFailedError()
	}

	if *invalidTokenCheck == false {
		return utils.InvalidTokenError()
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, utils.TokenSignatureError()
		}

		return []byte(os.Getenv("signing_secret")), nil
	})

	if err != nil {
		return utils.ParseTokenError(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["Id"] == userID {
			return nil
		}
	}

	return utils.InvalidTokenError()
}

func getTokenFromAuthHeader(authHeader string) (string, error) {
	var tokenString string
	headerValues := strings.Split(authHeader, " ")
	if len(headerValues) > 0 &&
		strings.ToLower(headerValues[0]) == "bearer" {
		tokenString = headerValues[1]
	} else {
		return "", utils.AuthHeaderError()
	}
	return tokenString, nil
}

func isInvalidToken(token string) *bool {
	db := Init()
	var invalidToken InvalidToken

	result := db.Table("invalid_tokens").Where("token = ?", token).Find(&invalidToken)
	res := result.RowsAffected == 0
	return &res
}
