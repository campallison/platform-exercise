package platform_exercise

import (
	"regexp"

	"github.com/campallison/platform-exercise/utils"
	"golang.org/x/crypto/bcrypt"
)

const (
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
	res, _ := regexp.MatchString(regex, name)
	return res
}

func ValidateEmail(req ValidateEmailRequest) (bool, error) {
	parsedEmail, err := utils.ParseEmail(req.Email)
	if err != nil {
		return false, utils.CouldNotParseEmailError(req.Email)
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

func GetUser(req GetUserRequest) (User, error) {
	db := Init()
	var user User

	if req.ID != "" {
		if err := db.Table("users").Where("id = ?", req.ID).First(&user).Error; err != nil {
			return User{}, utils.UserNotFoundError(req.ID)
		}
	}

	return user, nil
}

func UpdateUser(req UpdateUserRequest) (User, error) {
	db := Init()

	if req.ID != "" &&
		req.Name == "" &&
		req.Email == "" &&
		req.NewPassword == "" &&
		req.OldPassword == "" {
		return User{}, nil
	}

	var existing User
	if err := db.Where(`id = ?`, req.ID).First(&existing).Error; err != nil {
		return existing, utils.UserNotFoundError(req.ID)
	}

	if req.Name != "" && !isValidName(req.Name) {
		return User{}, utils.InvalidNameError(req.Name)
	}

	var hashedPW string
	if req.NewPassword != "" {
		if utils.PasswordStrength(req.NewPassword) < insecurePasswordThreshold {
			return User{}, utils.InsecurePasswordError()
		} else {
			err := bcrypt.CompareHashAndPassword([]byte(existing.Password), []byte(req.OldPassword))
			if err != nil {
				return User{}, utils.UnauthorizedError()
			} else {
				hashedPW, err = HashPassword(req.NewPassword)
				if err != nil {
					return User{}, err
				}
			}
		}
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

	if hashedPW != "" {
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

func DeleteUser(req DeleteUserRequest) (string, error) {
	db := Init()
	var user User
	if err := db.Where(`id = ?`, req.ID).First(&user).Error; err != nil {
		return req.ID, utils.UserNotFoundError(req.ID)
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := db.Delete(&User{}, "id = ?", req.ID).Error; err != nil {
		tx.Rollback()
		return req.ID, err
	}

	if err := tx.Commit().Error; err != nil {
		return req.ID, err
	}

	db.Unscoped().Where("id = ?", user.ID).First(&user)
	if user.ID == "" || !user.DeletedAt.Valid {
		tx.Rollback()
		return req.ID, utils.UserNotFoundError(user.ID)
	}

	return req.ID, nil
}
