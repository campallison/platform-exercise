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

func ValidateEmail(req ValidateEmailRequest) (bool, error) {
	parsedEmail, err := utils.ParseEmail(req.Email)
	if err != nil {
		return false, CouldNotParseEmailError(req.Email)
	}

	if utils.IsAliasedEmail(parsedEmail.LocalPart) {
		return false, AliasedEmailError(req.Email)
	}

	if utils.IsKnownSpamEmail(parsedEmail) {
		return false, ProhibitedEmailError(req.Email)
	}

	return true, nil
}

func CreateUser(db *gorm.DB, req CreateUserRequest) (User, error) {
	var user User

	if !isValidName(req.Name) {
		return User{}, InvalidNameError(req.Name)
	}

	user.Name = req.Name

	isValidEmail, err := ValidateEmail(ValidateEmailRequest{Email: req.Email})
	if err != nil {
		return User{}, err
	}

	if isValidEmail {
		user.Email = req.Email
	}

	if utils.PasswordStrength(req.Password) < insecurePasswordThreshold {
		return User{}, InsecurePasswordError()
	}

	hashedPW, err := HashPassword(req.Password)
	if err != nil {
		return User{}, err
	}

	user.Password = hashedPW

	if err := db.Save(&user).Error; err != nil {
		return User{}, SaveUserToDBError(user.Email)
	}

	return user, nil
}

func UpdateUser(db *gorm.DB, req UpdateUserRequest) (User, error) {
	// TODO implement token check in here
	var existing User
	if err := db.Where(`id = ?`, req.ID).First(&existing).Error; err != nil {
		return existing, UserNotFoundError(req.ID)
	}

	if req.ID != "" &&
		req.Name == "" &&
		req.Email == "" &&
		req.Password == "" {
		return User{}, nil
	}

	if req.Name != "" && !isValidName(req.Name) {
		return User{}, InvalidNameError(req.Name)
	}

	if req.Password != "" &&
		utils.PasswordStrength(req.Password) < insecurePasswordThreshold {
		return User{}, InsecurePasswordError()
	}

	hashedPW, err := HashPassword(req.Password)
	if err != nil {
		return User{}, err
	}

	if req.Email != "" {
		if ok, err := ValidateEmail(ValidateEmailRequest{Email: req.Email}); !ok {
			return User{}, err
		}
	}

	fields := map[string]interface{}{}

	if req.Name != "" {
		fields["name"] = req.Name
	}

	if req.Password != "" {
		fields["password"] = hashedPW
	}

	if req.Email != "" {
		fields["email"] = req.Email
	}

	db.Model(&existing).Where(`id = ?`, req.ID).Updates(fields)

	var updated User
	db.Table("users").Where("id = ?", req.ID).First(&updated)

	return updated, nil
}

func DeleteUser(db *gorm.DB, req DeleteUserRequest) (User, error) {
	return User{}, nil
}
