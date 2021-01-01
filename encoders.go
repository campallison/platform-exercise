package platform_exercise

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gt=0"`
}

type CreateUserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserRequest struct {
	ID       string `json:"id" validate:"required"`
	Name     string `json:"name"`
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"gt=0"`
}

type UpdateUserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type DeleteUserRequest struct {
	ID string `json:"id" validate:"required"`
}

type DeleteUserResponse struct {
	ID string `json:"id"`
}

type ValidateEmailRequest struct {
	Email string `json:"email" validate:"email"`
}

type ValidateEmailResponse struct {
	Email   string `json:"email"`
	IsValid bool   `json:"isValid`
	Error   string `json:"error"`
}

type PasswordStrengthRequest struct {
	Password string `json:"password" validate:"required"`
}

type PasswordStrengthResponse struct {
	Strength int `json:"strength"`
}
