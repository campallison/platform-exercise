package platform_exercise

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/campallison/platform-exercise/utils"
	"github.com/google/go-cmp/cmp"
	"gorm.io/gorm"
)

func Test_CreateUserHandler(t *testing.T) {
	databaseTest(t, func(database *gorm.DB) {
		clearDatabase(database)

		cases := []struct {
			name    string
			request events.APIGatewayProxyRequest
			status  int
			err     error
		}{
			{
				name: "successful call",
				request: events.APIGatewayProxyRequest{
					HTTPMethod: "POST",
					Headers:    map[string]string{"Content-Type": "application/json"},
					Body:       `{"name": "Leo Fender", "email": "leo@fender.com", "password": "nocastermaster46"}`,
				},
				status: 200,
			},
			{
				name: "bad email returns 422",
				request: events.APIGatewayProxyRequest{
					HTTPMethod: "POST",
					Headers:    map[string]string{"Content-Type": "application/json"},
					Body:       `{"name": "Leo Fender", "email": "leo@fender", "password": "nocastermaster46"}`,
				},
				status: 422,
			},
			{
				name: "invalid name returns 400",
				request: events.APIGatewayProxyRequest{
					HTTPMethod: "POST",
					Headers:    map[string]string{"Content-Type": "application/json"},
					Body:       `{"name": "Leo Fender)(*&", "email": "leo@fender.com", "password": "nocastermaster46"}`,
				},
				status: 400,
			},
			{
				name: "weak password returns 400",
				request: events.APIGatewayProxyRequest{
					HTTPMethod: "POST",
					Headers:    map[string]string{"Content-Type": "application/json"},
					Body:       `{"name": "Leo Fender", "email": "leo@fender.com", "password": "f"}`,
				},
				status: 400,
			},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				clearDatabase(database)

				response, _ := CreateUserHandler(c.request)

				if diff := cmp.Diff(c.status, response.StatusCode); diff != "" {
					t.Errorf("\nunexpected response (-want, +got)\n%s", diff)
				}
			})
		}
	})
}

func Test_ValidateEmailHandler(t *testing.T) {
	cases := []struct {
		name     string
		request  events.APIGatewayProxyRequest
		expected events.APIGatewayProxyResponse
		err      error
	}{
		{
			name: "successful call",
			request: events.APIGatewayProxyRequest{
				HTTPMethod: "POST",
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       `{"email":"leo@fender.com"}`,
			},
			expected: events.APIGatewayProxyResponse{
				Body:       `{"email":"leo@fender.com","isValid":true,"error":""}`,
				StatusCode: 200,
			},
			err: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			response, err := ValidateEmailHandler(c.request)

			utils.AssertErrorsEqual(t, c.err, err)

			if diff := cmp.Diff(c.expected, response); diff != "" {
				t.Errorf("\nunexpected response (-want, +got)\n%s", diff)
			}
		})
	}
}

func Test_PasswordStrengthHandler(t *testing.T) {
	cases := []struct {
		name     string
		request  events.APIGatewayProxyRequest
		expected events.APIGatewayProxyResponse
		err      error
	}{
		{
			name: "successful call",
			request: events.APIGatewayProxyRequest{
				HTTPMethod: "POST",
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       `{"password": "ArbitraryPassw0rd2Check!"}`,
			},
			expected: events.APIGatewayProxyResponse{
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       `{"strength":4}`,
				StatusCode: 200,
			},
			err: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			response, err := PasswordStrengthHandler(c.request)

			utils.AssertErrorsEqual(t, c.err, err)

			if diff := cmp.Diff(c.expected, response); diff != "" {
				t.Errorf("\nunexpected response (-want, +got)\n%s", diff)
			}
		})
	}
}
