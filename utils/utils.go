package utils

import (
	"regexp"
	"strings"

	"github.com/trustelem/zxcvbn"
	"golang.org/x/crypto/bcrypt"
)

var (
	aliasRegexp            = regexp.MustCompile(`\+`)
	emailRegexp            = regexp.MustCompile("(?i)([A-Z0-9._%+-]+@[A-Z0-9.-]+\\.[A-Z]{2,24})")
	prohibitedEmailDomains = []string{
		"0box.eu",
		"10minutemail.com",
		"anonbox.net",
		"contbay.com",
		"damnthespam.com",
		"dispostable.com",
		"fakemailgenerator.com",
		"grr.la",
		"guerillamail.biz",
		"guerillamail.com",
		"guerillamail.de",
		"guerillamail.info",
		"guerillamail.net",
		"guerillamail.org",
		"guerillamailblock.com",
		"hjdosage.com",
		"koszmail.pl",
		"kurzepost.de",
		"mailcatch.com",
		"mailforspam.com",
		"mailinator.com",
		"objectmail.com",
		"pokemail.net",
		"proxymail.eu",
		"rcpt.at",
		"sharklasers.com",
		"spamavert.com",
		"spam4.me",
		"trash-mail.at",
		"trash-mail.com",
		"trashmail.com",
		"trashmail.io",
		"trashmail.me",
		"trashmail.net",
		"trbvn.com",
		"urhen.com",
		"wegwerfmail.de",
		"wegwerfmail.net",
		"wegwerfmail.org",
		"yopmail.com",
	}
)

type Email struct {
	LocalPart string
	Domain    string
}

func IsAliasedEmail(email string) bool {
	return aliasRegexp.Match([]byte(email))
}

func PasswordStrength(password string) int {
	return zxcvbn.PasswordStrength(password, nil).Score
}

func IsKnownSpamEmail(email Email) bool {
	for _, domain := range prohibitedEmailDomains {
		if email.Domain == domain {
			return true
		}
	}

	return false
}

func ParseEmail(email string) (Email, error) {

	if !emailRegexp.MatchString(email) {
		return Email{}, InvalidEmailError(email)
	}

	if IsAliasedEmail(email) {
		return Email{}, AliasedEmailError(email)
	}

	i := strings.LastIndex(email, "@")
	return Email{LocalPart: email[:i], Domain: email[i+1:]}, nil
}

const (
	bcryptGenerationCost int = 14
)

type Credential struct {
	Hash   string `json:"-"`
	UserID string `json:"user_id"`
}

func (c Credential) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(c.Hash), []byte(password))
	return err == nil
}

func (c *Credential) HashPassword(password string) error {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcryptGenerationCost)
	if err != nil {
		return err
	}
	c.Hash = string(b)
	return nil
}
