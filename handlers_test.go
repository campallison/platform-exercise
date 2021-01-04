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

func Test_GetUserHandler(t *testing.T) {
	databaseTest(t, func(database *gorm.DB) {
		clearDatabase(database)

		email := "leo@fender.com"
		password := "SkunkStripeMapleNeckRosewoodFingerboard"
		user := User{
			Name:     "Leo Fender",
			Email:    email,
			Password: password,
		}
		database.Save(&user)

		token := utils.CreateTestToken(user.ID, email)

		aj := "application/json"
		validTokenHeader := utils.CreateTestAuthHeader(token, aj)

		cases := []struct {
			name    string
			request events.APIGatewayProxyRequest
			headers map[string]string
			status  int
		}{
			{
				name: "successful call with valid token",
				request: events.APIGatewayProxyRequest{
					HTTPMethod:     "GET",
					Headers:        validTokenHeader,
					PathParameters: map[string]string{"id": user.ID},
				},
				headers: validTokenHeader,
				status:  200,
			},
			{
				name: "error for call with valid token",
				request: events.APIGatewayProxyRequest{
					HTTPMethod:     "GET",
					Headers:        utils.CreateTestAuthHeader("invalidtoken", aj),
					PathParameters: map[string]string{"id": user.ID},
				},
				headers: validTokenHeader,
				status:  401,
			},
			{
				name: "error for user ID not found",
				request: events.APIGatewayProxyRequest{
					HTTPMethod:     "GET",
					Headers:        validTokenHeader,
					PathParameters: map[string]string{"id": "8b8b2419-0633-47fb-8f0f-7a515f2ccaa1"},
				},
				headers: validTokenHeader,
				status:  401,
			},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				response, _ := GetUserHandler(c.request)

				if diff := cmp.Diff(c.status, response.StatusCode); diff != "" {
					t.Errorf("\nunexpected response (-want, +got)\n%s", diff)
				}
			})
		}
	})
}

func Test_UpdateUserHandler(t *testing.T) {
	databaseTest(t, func(database *gorm.DB) {
		clearDatabase(database)

		email := "leo@fender.com"
		password := "SkunkStripeMapleNeckRosewoodFingerboard"
		user := User{
			Name:     "Leo Fender",
			Email:    email,
			Password: password,
		}
		database.Save(&user)

		token := utils.CreateTestToken(user.ID, email)

		aj := "application/json"
		validTokenHeader := utils.CreateTestAuthHeader(token, aj)

		cases := []struct {
			name    string
			request events.APIGatewayProxyRequest
			headers map[string]string
			status  int
		}{
			{
				name: "successful call with valid token",
				request: events.APIGatewayProxyRequest{
					HTTPMethod:     "PATCH",
					Headers:        validTokenHeader,
					PathParameters: map[string]string{"id": user.ID},
					Body:           `{"name": "Clarence Fender"}`,
				},
				headers: validTokenHeader,
				status:  200,
			},
			{
				name: "error for call with valid token",
				request: events.APIGatewayProxyRequest{
					HTTPMethod:     "PATCH",
					Headers:        utils.CreateTestAuthHeader("invalidtoken", aj),
					PathParameters: map[string]string{"id": user.ID},
					Body:           `{"name": "Clarence L. Fender"}`,
				},
				headers: validTokenHeader,
				status:  401,
			},
			{
				name: "error for user ID not found",
				request: events.APIGatewayProxyRequest{
					HTTPMethod:     "PATCH",
					Headers:        validTokenHeader,
					PathParameters: map[string]string{"id": "8b8b2419-0633-47fb-8f0f-7a515f2ccaa1"},
					Body:           `{"name": "Clarence L. Fender"}`,
				},
				headers: validTokenHeader,
				status:  401,
			},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				response, _ := UpdateUserHandler(c.request)

				if diff := cmp.Diff(c.status, response.StatusCode); diff != "" {
					t.Errorf("\nunexpected response (-want, +got)\n%s", diff)
				}
			})
		}
	})
}

func Test_DeleteUserHandler(t *testing.T) {
	databaseTest(t, func(database *gorm.DB) {
		clearDatabase(database)

		email := "leo@fender.com"
		password := "SkunkStripeMapleNeckRosewoodFingerboard"
		user := User{
			Name:     "Leo Fender",
			Email:    email,
			Password: password,
		}
		database.Save(&user)
		token := utils.CreateTestToken(user.ID, email)
		aj := "application/json"
		validTokenHeader := utils.CreateTestAuthHeader(token, aj)

		user2 := User{
			Name:     "Temp Dude",
			Email:    "tempdude@gmail.com",
			Password: "SoISaysToMableISays1234",
		}
		database.Save(&user2)
		token2 := utils.CreateTestToken(user2.ID, user2.Email)
		user2TokenHeader := utils.CreateTestAuthHeader(token2, aj)

		cases := []struct {
			name    string
			request events.APIGatewayProxyRequest
			headers map[string]string
			status  int
		}{
			{
				name: "successful call with valid token",
				request: events.APIGatewayProxyRequest{
					HTTPMethod:     "GET",
					Headers:        validTokenHeader,
					PathParameters: map[string]string{"id": user.ID},
				},
				headers: validTokenHeader,
				status:  200,
			},
			{
				name: "error for call with valid token",
				request: events.APIGatewayProxyRequest{
					HTTPMethod:     "GET",
					Headers:        utils.CreateTestAuthHeader("invalidtoken", aj),
					PathParameters: map[string]string{"id": user2.ID},
					Body:           `{"name": "Clarence L. Fender"}`,
				},
				headers: validTokenHeader,
				status:  401,
			},
			{
				name: "error for user ID not found",
				request: events.APIGatewayProxyRequest{
					HTTPMethod:     "GET",
					Headers:        user2TokenHeader,
					PathParameters: map[string]string{"id": "8b8b2419-0633-47fb-8f0f-7a515f2ccaa1"},
					Body:           `{"name": "Clarence L. Fender"}`,
				},
				headers: validTokenHeader,
				status:  401,
			},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				response, _ := DeleteUserHandler(c.request)

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
