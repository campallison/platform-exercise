package platform_exercise

import (
	"regexp"

	"github.com/campallison/platform-exercise/utils"
	"golang.org/x/crypto/bcrypt"
)

const (
	// zxcvbn threshold goes from:
	//   0: too guessable (guesses < 10^6)
	//   4: very unguessable (guesses >= 10^10)
	insecurePasswordThreshold = 2
	bcryptGenerationCost      = 14
)

func CheckPasswordStrength(password string) (err error) {
	if utils.PasswordStrength(password) < insecurePasswordThreshold {
		err = utils.InsecurePasswordError()
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
	// credit: https://stackoverflow.com/questions/2385701/regular-expression-for-first-and-last-name#comment103416432_45871742
	res, _ := regexp.MatchString(regex, name)
	return res
}

func ValidateEmail(req ValidateEmailRequest) (bool, error) {
	parsedEmail, err := utils.ParseEmail(req.Email)
	if err != nil {
		return false, utils.CouldNotParseEmailError(req.Email, err)
	}

	if utils.IsAliasedEmail(parsedEmail.LocalPart) {
		return false, utils.AliasedEmailError(req.Email)
	}

	if utils.IsKnownSpamEmail(parsedEmail) {
		return false, utils.ProhibitedEmailError(req.Email)
	}

	return true, nil
}

func CreateUser(req CreateUserRequest) (User, error) {
	db := Init()
	var user User

	if !isValidName(req.Name) {
		return User{}, utils.InvalidNameError(req.Name)
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
		return User{}, utils.InsecurePasswordError()
	}

	hashedPW, err := HashPassword(req.Password)
	if err != nil {
		return User{}, err
	}

	user.Password = hashedPW

	if err := db.Save(&user).Error; err != nil {
		return User{}, utils.SaveUserToDBError(user.Email)
	}

	return user, nil
}

func UpdateUser(req UpdateUserRequest) (User, error) {
	db := Init()

	if req.ID != "" &&
		req.Name == "" &&
		req.Email == "" &&
		req.Password == "" {
		return User{}, nil
	}

	var existing User
	if err := db.Where(`id = ?`, req.ID).First(&existing).Error; err != nil {
		return existing, utils.UserNotFoundError(req.ID)
	}

	if req.Name != "" && !isValidName(req.Name) {
		return User{}, utils.InvalidNameError(req.Name)
	}

	if req.Password != "" &&
		utils.PasswordStrength(req.Password) < insecurePasswordThreshold {
		return User{}, utils.InsecurePasswordError()
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

func DeleteUser(req DeleteUserRequest) (User, error) {
	db := Init()
	var user User
	if err := db.Where(`id = ?`, req.ID).First(&user).Error; err != nil {
		return user, utils.UserNotFoundError(req.ID)
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Table("users").Where("id  = ?", user.ID).Delete(User{}).Error; err != nil {
		tx.Rollback()
		return User{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return User{}, err
	}

	db.Unscoped().Where("id = ?", user.ID).First(&user)
	if user.ID == "" || !user.DeletedAt.Valid {
		tx.Rollback()
		return user, utils.UserNotFoundError(user.ID)
	}

	return user, nil
}
