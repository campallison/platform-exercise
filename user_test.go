package platform_exercise

import (
	"os"
	"testing"

	"github.com/campallison/platform-exercise/utils"
	"github.com/google/go-cmp/cmp"
	"gorm.io/gorm"
)

var runDBTests bool
var postgresURL string

func init() {
	// DB_TEST can be set arbitrarily for testing. I usually set it to 1.
	if os.Getenv("DB_TEST") == "" {
		// If env is not set - skip init
		return
	}

	if os.Getenv("DB_TEST") != "" {
		runDBTests = true
		postgresURL = PostgresURL
	}
}

func databaseTest(t *testing.T, handler func(database *gorm.DB)) {
	if !runDBTests {
		t.Skip("skipping database test, pass DB_TEST env to run")
	}
	database := Init(postgresURL)
	handler(database)
}

func setup(t *testing.T, database *gorm.DB) {
	if !runDBTests {
		t.Skip("skipping database test, pass DB_TEST env to run")
	}

	clearDatabase(database)
}

func clearDatabase(database *gorm.DB) {
	session := database.Session(&gorm.Session{AllowGlobalUpdate: true})
	session.Unscoped().Delete(User{})
}

func Test_CreateUser(t *testing.T) {
	databaseTest(t, func(database *gorm.DB) {
		setup(t, database)
		strongPW := "s3tIt0nF!re&Play1tWithYourT33th"

		cases := []struct {
			name     string
			setup    func(*gorm.DB)
			req      CreateUserRequest
			expected User
			err      error
		}{
			{
				name: "returns an error if email cannot be parsed",
				req: CreateUserRequest{
					Name:     "Leo Fender",
					Email:    "leo@fender",
					Password: strongPW,
				},
				expected: User{},
				err:      CouldNotParseEmailError("leo@fender"),
			},
			{
				name: "returns an error if email is aliased",
				req: CreateUserRequest{
					Name:     "Leo Fender",
					Email:    "leo+tune@fender.com",
					Password: strongPW,
				},
				expected: User{},
				err:      CouldNotParseEmailError("leo+tune@fender.com"),
			},
			{
				name: "returns an error if email is prohibited",
				req: CreateUserRequest{
					Name:     "Prohibited Guy",
					Email:    "nuge@trashmail.com",
					Password: strongPW,
				},
				expected: User{},
				err:      ProhibitedEmailError("nuge@trashmail.com"),
			},
			{
				name: "returns an error if password is not strong enough",
				req: CreateUserRequest{
					Name:     "Johnny NoSecurity",
					Email:    "jhonny@aol.com",
					Password: "1234",
				},
				expected: User{},
				err:      InsecurePasswordError(),
			},
			{
				name: "returns an error if name is invalid",
				req: CreateUserRequest{
					Name:     "I am the greetest!",
					Email:    "brain@infosphere.net",
					Password: "NowAmLeavingEarthForeverForNoRaisin",
				},
				expected: User{},
				err:      InvalidNameError("I am the greetest!"),
			},
			{
				name: "returns an error if save to db fails",
				setup: func(db *gorm.DB) {
					db.Save(&User{
						Name:  "Jimi Hendrix",
						Email: "voodoochild@fire.com",
					})
				},
				req: CreateUserRequest{
					Name:     "Other Name",
					Email:    "voodoochild@fire.com",
					Password: strongPW,
				},
				expected: User{},
				err:      SaveUserToDBError("voodoochild@fire.com"),
			},
			{
				name: "successfully saves a user",
				req: CreateUserRequest{
					Name:     "Leo Fender",
					Email:    "leo@fender.com",
					Password: strongPW,
				},
				expected: User{
					Name:  "Leo Fender",
					Email: "leo@fender.com",
				},
				err: nil,
			},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				setup(t, database)
				if c.setup != nil {
					c.setup(database)
				}

				res, err := CreateUser(database, c.req)

				utils.AssertErrorsEqual(t, c.err, err)

				if diff := cmp.Diff(
					c.expected.Name,
					res.Name,
				); diff != "" {
					t.Errorf("\nUnexpected user (-want, +got)\n%s", diff)
				}
			})
		}
	})
}

func Test_checkPasswordStrength(t *testing.T) {
	cases := []struct {
		name  string
		input string
		err   error
	}{
		{
			name:  "weak password returns an error",
			input: "weak",
			err:   InsecurePasswordError(),
		},
		{
			name:  "acceptable password returns no error",
			input: "s3tIt0nF!re&Play1tWithYourT33th",
			err:   nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := checkPasswordStrength(c.input)

			utils.AssertErrorsEqual(t, c.err, res)
		})
	}
}

func Test_isValidName(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid name returns true",
			input:    "Leo Fender",
			expected: true,
		},
		{
			name:     "valid name with foreign characters returns true",
			input:    "陳大文",
			expected: true,
		},
		{
			name:     "valid name with foreign characters returns true 2",
			input:    "আবাসযোগ্য",
			expected: true,
		},
		{
			name:     "valid name with foreign characters returns true 3",
			input:    "Biréli Lagrène",
			expected: true,
		},
		{
			name:     "invalid name returns false",
			input:    "A$@p Rocky",
			expected: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := isValidName(c.input)

			if diff := cmp.Diff(c.expected, res); diff != "" {
				t.Errorf("\nUnexpected result (-want, +got)\n%s", diff)
			}
		})
	}
}
