package typ

import (
	"strings"

	"github.com/dusted-go/fault/fault"
)

var (
	EmptyEmail = Email{
		value:  "",
		domain: "",
	}
)

// Email defines an email address object.
type Email struct {
	value  string
	domain string
}

// Normalised returns a lowercase and trimmed string representation of an email address.
func (e Email) Normalised() string {
	return e.value
}

// Domain returns the part of an email address after the '@' sign.
func (e Email) Domain() string {
	return e.domain
}

// Equals checks if two email strings are the same.
func (e Email) Equals(other string) bool {
	return e.value == strings.ToLower(other)
}

// NewEmail validates, normalises and creates a new email.
func NewEmail(value string) (Email, error) {

	invalidEmailFault := fault.User("invalid_email_address", "Email address is invalid.")
	value = strings.TrimSpace(strings.ToLower(value))
	length := len(value)

	if length == 0 {
		return EmptyEmail,
			fault.User("missing_email_address", "Email address is required.")
	}

	// Assuming minimum email is: x@x.xx
	if length < 6 {
		return EmptyEmail, invalidEmailFault
	}

	if !strings.ContainsRune(value, '@') {
		return EmptyEmail, invalidEmailFault
	}

	if !strings.ContainsRune(value, '.') {
		return EmptyEmail, invalidEmailFault
	}

	if strings.LastIndex(value, ".") < strings.LastIndex(value, "@") {
		return EmptyEmail, invalidEmailFault
	}

	domain := strings.SplitN(value, "@", 2)[1]
	email := Email{
		value:  value,
		domain: domain,
	}

	return email, nil
}
