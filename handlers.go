package platform_exercise

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/campallison/platform-exercise/utils"
)

func badRequestResponse(err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Body:       err.Error(),
	}, nil
}

func unauthorizedResponse() (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusUnauthorized,
	}, nil
}

func CreateUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var createUserReq CreateUserRequest
	if err := json.Unmarshal([]byte(request.Body), &createUserReq); err != nil {
		return badRequestResponse(err)
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

func GetUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var getUserReq GetUserRequest
	getUserReq.ID = request.PathParameters["id"]

	if err := CheckToken(request.Headers["Authorization"], getUserReq.ID); err != nil {
		return unauthorizedResponse()
	}

	retrievedUser, err := GetUser(getUserReq)
	if err != nil {
		apiError := err.(utils.APIError)

		return events.APIGatewayProxyResponse{
			StatusCode: apiError.Code,
			Headers:    map[string]string{"Content-Type": "text/plain"},
			Body:       apiError.Message,
		}, nil
	}

	body, err := json.Marshal(GetUserResponse{
		ID:    retrievedUser.ID,
		Name:  retrievedUser.Name,
		Email: retrievedUser.Email,
	})

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(body),
	}, nil
}

func UpdateUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var updateUserReq UpdateUserRequest
	if err := json.Unmarshal([]byte(request.Body), &updateUserReq); err != nil {
		return badRequestResponse(err)
	}
	updateUserReq.ID = request.PathParameters["id"]

	if err := CheckToken(request.Headers["Authorization"], updateUserReq.ID); err != nil {
		return unauthorizedResponse()
	}

	updatedUser, err := UpdateUser(updateUserReq)
	if err != nil {
		apiError := err.(utils.APIError)

		return events.APIGatewayProxyResponse{
			StatusCode: apiError.Code,
			Headers:    map[string]string{"Content-Type": "text/plain"},
			Body:       apiError.Message,
		}, nil
	}

	body, err := json.Marshal(UpdateUserResponse{
		ID:    updatedUser.ID,
		Name:  updatedUser.Name,
		Email: updatedUser.Email,
	})

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(body),
	}, nil
}

func DeleteUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	deleteUserReq := DeleteUserRequest{ID: request.PathParameters["id"]}

	if err := CheckToken(request.Headers["Authorization"], deleteUserReq.ID); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
		}, nil
	}

	deletedUser, err := DeleteUser(deleteUserReq)
	if err != nil {
		apiError := err.(utils.APIError)

		return events.APIGatewayProxyResponse{
			StatusCode: apiError.Code,
			Headers:    map[string]string{"Content-Type": "text/plain"},
			Body:       apiError.Message,
		}, nil
	}

	body, err := json.Marshal(DeleteUserResponse{ID: deletedUser})

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(body),
	}, nil
}

func ValidateEmailHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var validateEmailReq ValidateEmailRequest
	err := json.Unmarshal([]byte(request.Body), &validateEmailReq)
	if err != nil {
		return badRequestResponse(err)
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
	pwStrength := utils.PasswordStrength(request.Body)

	body, _ := json.Marshal(PasswordStrengthResponse{Strength: pwStrength})

	return events.APIGatewayProxyResponse{
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(body),
		StatusCode: 200,
	}, nil
}

func LoginHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var creds Credential
	if err := json.Unmarshal([]byte(request.Body), &creds); err != nil {
		return badRequestResponse(err)
	}

	loginResult, err := Login(creds)
	if err != nil {
		return badRequestResponse(err)
	}

	body, _ := json.Marshal(loginResult)

	return events.APIGatewayProxyResponse{
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(body),
		StatusCode: 200,
	}, nil
}

func LogoutHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var logoutRequest LogoutRequest
	id := request.PathParameters["id"]
	authHeader := request.Headers["Authorization"]
	if err := CheckToken(authHeader, id); err != nil {
		return unauthorizedResponse()
	}

	logoutRequest.AccessToken, _ = getTokenFromAuthHeader(authHeader)
	logoutRequest.ID = id

	logoutResult, err := Logout(logoutRequest)
	if err != nil {
		return badRequestResponse(err)
	}

	body, _ := json.Marshal(logoutResult)

	return events.APIGatewayProxyResponse{
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(body),
		StatusCode: 200,
	}, nil
}
