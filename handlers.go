package platform_exercise

import (
	"github.com/aws/aws-lambda-go/events"
)

func CreateUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       `{"createUser": "hit"}`,
	}, nil
}

func ValidateEmailHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       `{"validateEmail": "hit"}`,
		StatusCode: 200,
	}, nil
}
