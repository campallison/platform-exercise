package platform_exercise

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/campallison/platform-exercise/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const PostgresURL = "postgres://root:postgres@localhost:5432/postgres?sslmode=disable"

const (
	// zxcvbn threshold goes from:
	//   0: too guessable (guesses < 10^6)
	//   4: very unguessable (guesses >= 10^10)
	insecurePasswordThreshold = 2
	bcryptGenerationCost      = 14
)

func CouldNotParseEmailError(email string) error {
	return errors.New(fmt.Sprintf("unable to parse email %s", email))
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
	return errors.New(fmt.Sprintf("Error saving user with email %s. Unknown reason, but may already exist.", email))
}

func checkPasswordStrength(password string) (err error) {
	if utils.PasswordStrength(password) < insecurePasswordThreshold {
		err = InsecurePasswordError()
	}
	return
}

func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcryptGenerationCost)
	if err != nil {
		return "", err
	}
	hash := string(b)
	return hash, nil
}

func isValidName(name string) bool {
	regex := `^[^0-9_!¡?÷?¿/\\+=@#$%ˆ&*(){}|~<>;:[\]]{2,}$`
	res, _ := regexp.MatchString(regex, name)
	return res
}

func CreateUser(db *gorm.DB, req CreateUserRequest) (User, error) {
	var user User

	if !isValidName(req.Name) {
		return User{}, InvalidNameError(req.Name)
	}

	parsedEmail, err := utils.ParseEmail(req.Email)
	if err != nil {
		return User{}, CouldNotParseEmailError(req.Email)
	}

	if utils.IsAliasedEmail(parsedEmail.LocalPart) {
		return User{}, AliasedEmailError(req.Email)
	}

	if utils.IsKnownSpamEmail(parsedEmail) {
		return User{}, ProhibitedEmailError(req.Email)
	}

	if utils.PasswordStrength(req.Password) < insecurePasswordThreshold {
		return User{}, InsecurePasswordError()
	}

	hashedPW, err := HashPassword(req.Password)
	if err != nil {
		return User{}, err
	}

	user.Name = req.Name
	user.Email = req.Email
	user.Password = hashedPW

	if err := db.Save(&user).Error; err != nil {
		return User{}, SaveUserToDBError(user.Email)
	}

	return user, nil
}
