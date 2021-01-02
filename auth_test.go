package platform_exercise

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_LoginTest(t *testing.T) {

}

func Test_CheckPassword(t *testing.T) {
	cases := []struct {
		name     string
		creds    Credential
		hash     string
		expected bool
	}{
		{
			name: "should return true for passwords that match",
			creds: Credential{
				Password: "ILoveMarge742!!!ScrewFlanders",
			},
			hash:     "$2a$14$dKA9DtzAy5dqcEh.Cd36RetU0okxCLsjww8WH8tuQKeB9AQNojdwy",
			expected: true,
		},
		{
			name: "should return false for passwords that do not match",
			creds: Credential{
				Password: "PoliticalSloganAndElectionYear!!",
			},
			hash:     "$2a$14$dKA9DtzAy5dqcEh.Cd36RetU0okxCLsjww8WH8tuQKeB9AQNojdwy",
			expected: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := c.creds.CheckPassword(c.hash)
			if diff := cmp.Diff(c.expected, result); diff != "" {
				t.Errorf("unexpected result (-want, +got)\n%s", diff)
			}
		})
	}
}
