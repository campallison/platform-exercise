package utils

import (
	errors2 "errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_passwordStrength(t *testing.T) {
	type args struct {
		password string
	}
	cases := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Test empty",
			args: args{""},
			want: 0,
		},
		{
			name: "Test password",
			args: args{"password"},
			want: 0,
		},
		{
			name: "Test guessable",
			args: args{"!keRp@"},
			want: 1,
		},
		{
			name: "Test somewhat guessable",
			args: args{"$h0rnSg4!"},
			want: 3,
		},
		{
			name: "Test secure",
			args: args{"this is pretty ok"},
			want: 4,
		},
		{
			name: "Test long without numbers and symbols",
			args: args{"correct battery horse stapler"},
			want: 4,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := PasswordStrength(c.args.password)
			if diff := cmp.Diff(c.want, res); diff != "" {
				t.Errorf("\nUnexpected password strength (-want, +got)\n%s", diff)
			}
		})
	}
}

func Test_isAliasedEmail(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "email contains an alias",
			input: "allison+fender@fender.com",
			want:  true,
		},
		{
			name:  "email does not contain an alias",
			input: "allison@fender.com",
			want:  false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := IsAliasedEmail(c.input)

			if diff := cmp.Diff(c.want, res); diff != "" {
				t.Errorf("\nUnexpected result (-want, +got)\n%s", diff)
			}
		})
	}
}

func Test_isKnownSpamEmail(t *testing.T) {
	cases := []struct {
		name  string
		input Email
		want  bool
	}{
		{
			name:  "doesn't consider good domain spam",
			input: Email{LocalPart: "leo", Domain: "fender.com"},
			want:  false,
		},
		{
			name:  "doesn't consider good domain spam",
			input: Email{LocalPart: "sergey", Domain: "google.com"},
			want:  false,
		},
		{
			name:  "doesn't consider good domain spam",
			input: Email{LocalPart: "leo", Domain: "tcell.io"},
			want:  false,
		},
		{
			name:  "detects known bad domain as spam",
			input: Email{LocalPart: "leo", Domain: "mailforspam.com"},
			want:  true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := IsKnownSpamEmail(c.input)

			if diff := cmp.Diff(c.want, res); diff != "" {
				t.Errorf("\nUnexpected result (-want, +got)\n%s", diff)
			}
		})
	}
}

func Test_parseEmail(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		err      error
		expected Email
	}{
		{
			name:     "valid fender email",
			input:    "leo@fender.com",
			err:      nil,
			expected: Email{LocalPart: "leo", Domain: "fender.com"},
		},
		{
			name:     "valid google email",
			input:    "larry@google.com",
			err:      nil,
			expected: Email{LocalPart: "larry", Domain: "google.com"},
		},
		{
			name:     "valid .io email",
			input:    "leo@tcell.io",
			err:      nil,
			expected: Email{LocalPart: "leo", Domain: "tcell.io"},
		},
		{
			name:  "no localpart",
			input: "@fender.com",
			err:   InvalidEmailError("@fender.com"),
		},
		{
			name:  "no domain",
			input: "leo@",
			err:   InvalidEmailError("leo@"),
		},
		{
			name:  "no SLD in domain",
			input: "leo@fender",
			err:   InvalidEmailError("leo@fender"),
		},
		{
			name:  "contains alias",
			input: "leo+fender@fender.com",
			err:   AliasedEmailError("leo+fender@fender.com"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res, err := ParseEmail(c.input)

			AssertErrorsEqual(t, c.err, err)

			if diff := cmp.Diff(c.expected, res); diff != "" {
				t.Errorf("\nUnexpected email (-want, +got)\n%s", diff)
			}
		})
	}
}
