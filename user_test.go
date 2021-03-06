package platform_exercise

import (
	"testing"

	"github.com/campallison/platform-exercise/utils"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gorm.io/gorm"
)

func databaseTest(t *testing.T, handler func(database *gorm.DB)) {
	database := Init()
	handler(database)
}

func clearDatabase(database *gorm.DB) {
	session := database.Session(&gorm.Session{AllowGlobalUpdate: true})
	session.Unscoped().Delete(User{})
}

func Test_CreateUser(t *testing.T) {
	databaseTest(t, func(database *gorm.DB) {
		clearDatabase(database)
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
				err:      utils.CouldNotParseEmailError("leo@fender"),
			},
			{
				name: "returns an error if email is aliased",
				req: CreateUserRequest{
					Name:     "Leo Fender",
					Email:    "leo+tune@fender.com",
					Password: strongPW,
				},
				expected: User{},
				err:      utils.CouldNotParseEmailError("leo+tune@fender.com"),
			},
			{
				name: "returns an error if email is prohibited",
				req: CreateUserRequest{
					Name:     "Prohibited Guy",
					Email:    "nuge@trashmail.com",
					Password: strongPW,
				},
				expected: User{},
				err:      utils.ProhibitedEmailError("nuge@trashmail.com"),
			},
			{
				name: "returns an error if password is not strong enough",
				req: CreateUserRequest{
					Name:     "Johnny NoSecurity",
					Email:    "jhonny@aol.com",
					Password: "1234",
				},
				expected: User{},
				err:      utils.InsecurePasswordError(),
			},
			{
				name: "returns an error if name is invalid",
				req: CreateUserRequest{
					Name:     "I am the greetest!",
					Email:    "brain@infosphere.net",
					Password: "NowAmLeavingEarthForeverForNoRaisin",
				},
				expected: User{},
				err:      utils.InvalidNameError("I am the greetest!"),
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
				err:      utils.SaveUserToDBError("voodoochild@fire.com"),
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
				clearDatabase(database)
				if c.setup != nil {
					c.setup(database)
				}

				res, err := CreateUser(c.req)
				utils.AssertErrorsEqual(t, c.err, err)
				if diff := cmp.Diff(
					c.expected,
					res,
					cmpopts.IgnoreFields(User{}, "ID"),
					cmpopts.IgnoreFields(User{}, "Password"),
					cmpopts.IgnoreFields(User{}, "CreatedAt"),
					cmpopts.IgnoreFields(User{}, "UpdatedAt"),
					cmpopts.IgnoreFields(User{}, "DeletedAt"),
				); diff != "" {
					t.Errorf("\nUnexpected user (-want, +got)\n%s", diff)
				}
			})
		}
	})
}

func Test_GetUser(t *testing.T) {
	databaseTest(t, func(database *gorm.DB) {
		clearDatabase(database)
		id := "5a135f2a-976d-4db7-9910-8a96b043b1b6"

		cases := []struct {
			name     string
			setup    func(*gorm.DB)
			req      GetUserRequest
			expected User
			err      error
		}{
			{
				name: "returns an error for user not found, by ID",
				req: GetUserRequest{
					ID: id,
				},
				expected: User{},
				err:      utils.UserNotFoundError(id),
			},
			{
				name: "retrieves a user by ID",
				setup: func(db *gorm.DB) {
					db.Save(&User{
						ID:       id,
						Name:     "Enrico Fermi",
						Email:    "ilpapa@umich.edu",
						Password: "ThePileIsCritical",
					})
				},
				req: GetUserRequest{
					ID: id,
				},
				expected: User{
					ID:    id,
					Name:  "Enrico Fermi",
					Email: "ilpapa@umich.edu",
				},
				err: nil,
			},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				clearDatabase(database)
				if c.setup != nil {
					c.setup(database)
				}

				res, err := GetUser(c.req)
				utils.AssertErrorsEqual(t, c.err, err)
				if diff := cmp.Diff(
					c.expected,
					res,
					cmpopts.IgnoreFields(User{}, "ID"),
					cmpopts.IgnoreFields(User{}, "Password"),
					cmpopts.IgnoreFields(User{}, "CreatedAt"),
					cmpopts.IgnoreFields(User{}, "UpdatedAt"),
					cmpopts.IgnoreFields(User{}, "DeletedAt"),
				); diff != "" {
					t.Errorf("\nUnexpected user (-want, +got)\n%s", diff)
				}
			})
		}
	})
}

func Test_UpdateUser(t *testing.T) {
	databaseTest(t, func(database *gorm.DB) {
		clearDatabase(database)
		id := "13a185dd-1c2e-4092-81cc-ec306d18b2bd"
		frysPW := "WalkinOnSunshine1999!"
		frysHash, _ := HashPassword(frysPW)

		cases := []struct {
			name     string
			setup    func(*gorm.DB)
			req      UpdateUserRequest
			expected User
			err      error
		}{
			{
				name: "should return an error if the requested user ID does not exist",
				req: UpdateUserRequest{
					ID:   id,
					Name: "Philip Fry",
				},
				expected: User{},
				err:      utils.UserNotFoundError(id),
			},
			{
				name: "should return empty user and no error if only ID is provided",
				setup: func(db *gorm.DB) {
					db.Save(&User{
						ID:       id,
						Name:     "Philip Fry",
						Email:    "deliveryboy@panuccis.net",
						Password: "WalkinOnSunshine1999",
					})
				},
				req: UpdateUserRequest{
					ID: id,
				},
				expected: User{},
				err:      nil,
			},
			{
				name: "does not allow invalid name",
				setup: func(db *gorm.DB) {
					db.Save(&User{
						ID:       id,
						Name:     "Philip Fry",
						Email:    "deliveryboy@panuccis.net",
						Password: "WalkinOnSunshine1999!",
					})
				},
				req: UpdateUserRequest{
					ID:   id,
					Name: "!nv@lid Name",
				},
				expected: User{},
				err:      utils.InvalidNameError("!nv@lid Name"),
			},
			{
				name: "does not allow new password if below strength threshold",
				setup: func(db *gorm.DB) {
					db.Save(&User{
						ID:       id,
						Name:     "Philip Fry",
						Email:    "deliveryboy@panuccis.net",
						Password: "WalkinOnSunshine1999!",
					})
				},
				req: UpdateUserRequest{
					ID:          id,
					OldPassword: "WalkinOnSunshine1999!",
					NewPassword: "weak",
				},
				expected: User{},
				err:      utils.InsecurePasswordError(),
			},
			{
				name: "returns an error if new email is invalid",
				setup: func(db *gorm.DB) {
					db.Save(&User{
						ID:       id,
						Name:     "Philip Fry",
						Email:    "deliveryboy@panuccis.net",
						Password: "WalkinOnSunshine1999!",
					})
				},
				req: UpdateUserRequest{
					ID:    id,
					Email: "bender@isGreat",
				},
				expected: User{},
				err:      utils.CouldNotParseEmailError("bender@isGreat"),
			},
			{
				name: "successfully updates a valid name",
				setup: func(db *gorm.DB) {
					db.Save(&User{
						ID:       id,
						Name:     "Philip Fry",
						Email:    "deliveryboy@panuccis.net",
						Password: "WalkinOnSunshine1999!",
					})
				},
				req: UpdateUserRequest{
					ID:   id,
					Name: "Bender Rodriguez",
				},
				expected: User{
					ID:    id,
					Name:  "Bender Rodriguez",
					Email: "deliveryboy@panuccis.net",
				},
				err: nil,
			},
			{
				name: "successfully updates a valid email",
				setup: func(db *gorm.DB) {
					db.Save(&User{
						ID:       id,
						Name:     "Philip Fry",
						Email:    "deliveryboy@panuccis.net",
						Password: "WalkinOnSunshine1999!",
					})
				},
				req: UpdateUserRequest{
					ID:    id,
					Email: "daffodil@shiny.com",
				},
				expected: User{
					ID:    id,
					Name:  "Philip Fry",
					Email: "daffodil@shiny.com",
				},
				err: nil,
			},
			{
				name: "successfully updates a valid name and email",
				setup: func(db *gorm.DB) {
					db.Save(&User{
						ID:       id,
						Name:     "Philip Fry",
						Email:    "deliveryboy@panuccis.net",
						Password: "WalkinOnSunshine1999!",
					})
				},
				req: UpdateUserRequest{
					ID:    id,
					Name:  "Bender Rodriguez",
					Email: "daffodil@shiny.com",
				},
				expected: User{
					ID:    id,
					Name:  "Bender Rodriguez",
					Email: "daffodil@shiny.com",
				},
				err: nil,
			},
			{
				name: "successfully updates a password given correct current password",
				setup: func(db *gorm.DB) {
					db.Save(&User{
						ID:       id,
						Name:     "Philip Fry",
						Email:    "deliveryboy@panuccis.net",
						Password: frysHash,
					})
				},
				req: UpdateUserRequest{
					ID:          id,
					OldPassword: "WalkinOnSunshine1999!",
					NewPassword: "BenderIsGreat3001!",
				},
				expected: User{
					ID:    id,
					Name:  "Philip Fry",
					Email: "deliveryboy@panuccis.net",
				},
				err: nil,
			},
			{
				name: "does not update a password given incorrect current password",
				setup: func(db *gorm.DB) {
					db.Save(&User{
						ID:       id,
						Name:     "Philip Fry",
						Email:    "deliveryboy@panuccis.net",
						Password: "$2a$14$8/kFdq99JkTajySO88HzS.PQntfZFt19FvEsshCYb3zG1NpR2/yiS",
					})
				},
				req: UpdateUserRequest{
					ID:          id,
					OldPassword: "WalkinOnSunshine1999!",
					NewPassword: "BenderIsGreat3001!",
				},
				expected: User{},
				err:      utils.UnauthorizedError(),
			},
			{
				name: "does not update a password given weak new password",
				setup: func(db *gorm.DB) {
					db.Save(&User{
						ID:       id,
						Name:     "Philip Fry",
						Email:    "deliveryboy@panuccis.net",
						Password: frysHash,
					})
				},
				req: UpdateUserRequest{
					ID:          id,
					OldPassword: "WalkinOnSunshine1999!",
					NewPassword: "TL",
				},
				expected: User{},
				err:      utils.InsecurePasswordError(),
			},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				clearDatabase(database)
				if c.setup != nil {
					c.setup(database)
				}

				res, err := UpdateUser(c.req)
				utils.AssertErrorsEqual(t, c.err, err)
				if diff := cmp.Diff(
					c.expected,
					res,
					cmpopts.IgnoreFields(User{}, "Password"),
					cmpopts.IgnoreFields(User{}, "CreatedAt"),
					cmpopts.IgnoreFields(User{}, "UpdatedAt"),
					cmpopts.IgnoreFields(User{}, "DeletedAt"),
				); diff != "" {
					t.Errorf("\nUnexpected user (-want, +got)\n%s", diff)
				}
			})
		}
	})
}

func Test_DeleteUser(t *testing.T) {
	databaseTest(t, func(database *gorm.DB) {
		clearDatabase(database)
		id := "13a185dd-1c2e-4092-81cc-ec306d18b2bd"

		cases := []struct {
			name  string
			setup func(*gorm.DB)
			req   DeleteUserRequest
			err   error
		}{
			{
				name: "returns an error if requested user ID is not found",
				req:  DeleteUserRequest{ID: id},
				err:  utils.UserNotFoundError(id),
			},
			{
				name: "deletes a user successfully",
				setup: func(db *gorm.DB) {
					db.Save(&User{
						ID:       id,
						Name:     "Ephemeral User",
						Email:    "user@domain.com",
						Password: "irrelevant",
					})
				},
				req: DeleteUserRequest{ID: id},
				err: nil,
			},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				clearDatabase(database)
				if c.setup != nil {
					c.setup(database)
				}

				_, err := DeleteUser(c.req)
				utils.AssertErrorsEqual(t, c.err, err)
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
			err:   utils.InsecurePasswordError(),
		},
		{
			name:  "acceptable password returns no error",
			input: "s3tIt0nF!re&Play1tWithYourT33th",
			err:   nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := CheckPasswordStrength(c.input)

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
