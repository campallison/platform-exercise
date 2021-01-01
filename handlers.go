package platform_exercise

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/campallison/platform-exercise/utils"
)

func CreateUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var createUserReq CreateUserRequest
	err := json.Unmarshal([]byte(request.Body), &createUserReq)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}

	createdUser, err := CreateUser(createUserReq)
	if err != nil {
		apiError := err.(utils.APIError)

		return events.APIGatewayProxyResponse{
			StatusCode: apiError.Code,
			Headers:    map[string]string{"Content-Type": "text/plain"},
			Body:       apiError.Message,
		}, nil
	}

	body, err := json.Marshal(CreateUserResponse{
		ID:    createdUser.ID,
		Name:  createdUser.Name,
		Email: createdUser.Email,
	})

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(body),
	}, nil
}

func UpdateUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       `{"updateUserHandler": "hit"}`,
		StatusCode: 200,
	}, nil
}

func DeleteUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       `{"deleteUserHandler": "hit"}`,
		StatusCode: 200,
	}, nil
}

func ValidateEmailHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var validateEmailReq ValidateEmailRequest
	err := json.Unmarshal([]byte(request.Body), &validateEmailReq)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}

	isEmailValid, err := ValidateEmail(validateEmailReq)
	errToRespond := ""
	if err != nil {
		errToRespond = err.Error()
	}
	body, err := json.Marshal(ValidateEmailResponse{
		Email:   validateEmailReq.Email,
		IsValid: isEmailValid,
		Error:   errToRespond,
	})

	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
	}, nil
}

func PasswordStrengthHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       `{"passwordStrengthHandler": "hit"}`,
		StatusCode: 200,
	}, nil
}
