package email

import (
	"strings"

	"github.com/dusted-go/utils/fault"
)

var (
	Empty = Address{
		value:  "",
		domain: "",
	}
)

// Address defines an email address object.
type Address struct {
	value  string
	domain string
}

// Normalised returns a lowercase and trimmed string representation of an email address.
func (a Address) Normalised() string {
	return a.value
}

// Domain returns the part of an email address after the '@' sign.
func (a Address) Domain() string {
	return a.domain
}

// Equals checks if two email strings are the same.
func (a Address) Equals(other string) bool {
	return a.value == strings.ToLower(other)
}

// New validates, normalises and creates a new address.
func New(value string) (Address, error) {

	invalidEmailFault := fault.User("invalid_email_address", "Email address is invalid.")
	value = strings.TrimSpace(strings.ToLower(value))
	length := len(value)

	if length == 0 {
		return Empty,
			fault.User("missing_email_address", "Email address is required.")
	}

	// Assuming minimum email is: x@x.xx
	if length < 6 {
		return Empty, invalidEmailFault
	}

	if !strings.ContainsRune(value, '@') {
		return Empty, invalidEmailFault
	}

	if !strings.ContainsRune(value, '.') {
		return Empty, invalidEmailFault
	}

	if strings.LastIndex(value, ".") < strings.LastIndex(value, "@") {
		return Empty, invalidEmailFault
	}

	domain := strings.SplitN(value, "@", 2)[1]
	emailAddr := Address{
		value:  value,
		domain: domain,
	}

	return emailAddr, nil
}
