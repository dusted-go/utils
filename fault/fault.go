package fault

import (
	"fmt"
	"io"
	"strings"

	pkgerrors "github.com/pkg/errors"
)

// ------
// User Error
// ------

// UserError represents an error that was typically caused by the end user.
// A user error is normally a type of error which an application would like to surface back
// to the user. It could be something like a validation error of some user provided input or
// other errors that would normally result in a 4xx status code in a web application context.
//
// User errors should also contain an error code in addition to the error message.
// This helps the parsing of user errors by external programs which act on behalf of the user.
// For example a UserError returned by a HTTP API would return an error code alongside the message
// so that the calling client can parse the error and decide what to do next.
//
// Error codes should ideally be unique and descriptive strings in order to prevent collision in a larger application.
//
// Examples:
//    "MISSING_FIRST_NAME": "Please provide your first name"
//    "INVALID_EMAIL_ADDR": "Please provide a valid email address"
type UserError struct {
	// map of error codes and messages
	errors map[string]string

	// codes is used to preserve the order in which
	// errors are being added, since a map[string]string
	// will iterate in random order.
	codes []string
}

// Add appends an additional user error to the collection of errors.
func (e *UserError) Add(code string, msg string) {
	e.codes = append(e.codes, code)
	e.errors[code] = msg
}

// Addf appends an additional user error to the collection of errors.
func (e *UserError) Addf(code string, format string, a ...interface{}) {
	e.Add(code, fmt.Sprintf(format, a...))
}

func (e *UserError) errorMessage(includeCode bool) string {
	if len(e.errors) == 0 {
		return ""
	}
	prefix := "- "
	if len(e.errors) == 1 {
		prefix = ""
	}
	sb := strings.Builder{}
	for _, k := range e.codes {
		v := e.errors[k]
		if sb.Len() > 0 {
			sb.WriteString("\n")
		}
		if includeCode {
			sb.WriteString(fmt.Sprintf("%s%s (%s)", prefix, v, k))
		} else {
			sb.WriteString(fmt.Sprintf("%s%s", prefix, v))
		}
	}
	return sb.String()
}

// String returns the error message.
func (e *UserError) String() string {
	return e.Error()
}

// Error will return a string of one or all user errors.
//
// If there is only one user error it will be represented as a single string.
//   Example:
//      Email address is required (MISSING_EMAIL_ADDRESS)
//
// If there are more than one user error (e.g. multiple validation errors)
// then a multi line string resembling a list of errors will be returned.
//   Example:
//      - First name is required (MISSING_FIRST_NAME)
//      - Last name is required (MISSING_LAST_NAME)
//      - Invalid email address (INVALID_EMAIL_ADDRESS)
//
// Use FriendlyError() to compute the same string without error codes attached.
//
// Use ErrorMessages() to get an array of the messages only (no codes attached).
func (e *UserError) Error() string {
	return e.errorMessage(true)
}

// FriendlyError will return a string of one or all user errors.
//
// FriendlyError is equivalent to Error() except it doesn't include error codes in the message.
//
// If there is only one user error it will be represented as a single string.
//   Example:
//      Email address is required
//
// If there are more than one user error (e.g. multiple validation errors)
// then a multi line string resembling a list of errors will be returned.
//   Example:
//      - First name is required
//      - Last name is required
//      - Invalid email address
//
// Use Error() to compute the same string with error codes attached.
//
// Use ErrorMessages() to get an array of the messages only (no codes attached).
func (e *UserError) FriendlyError() string {
	return e.errorMessage(false)
}

// Errors returns a map of error codes and messages.
func (e *UserError) Errors() map[string]string {
	return e.errors
}

// ErrorMessages returns an array of error messages only.
func (e *UserError) ErrorMessages() []string {
	messages := make([]string, len(e.codes))
	for i, k := range e.codes {
		messages[i] = e.errors[k]
	}
	return messages
}

// User creates a new UserError fault.
func User(code string, msg string) *UserError {
	return &UserError{
		errors: map[string]string{
			code: msg,
		},
		codes: []string{code},
	}
}

// Userf creates a new UserError fault.
func Userf(code string, format string, a ...interface{}) *UserError {
	return User(code, fmt.Sprintf(format, a...))

}

// ------
// System Error
// ------

const (
	padding = "   "
)

// SystemError represents an error that was caused by an internal fault.
// A system error is typically an error which can only be handled by the application
// itself or would typically result in a 5xx status code in a web application context.
//
// Examples:
// - error connecting to a database
// - error reading from an IO stream
// - unexpected error from making a HTTP call
// - etc.
type SystemError struct {
	err  error
	msgs []string
}

// Cause returns the underlying error.
func (e *SystemError) Cause() error {
	return e.err
}

// String returns the error message.
func (e *SystemError) String() string {
	return e.Error()
}

// Error returns the error message.
func (e *SystemError) Error() string {
	return e.err.Error()
}

// StackTrace returns the error message including the stack trace.
func (e *SystemError) StackTrace() string {
	return fmt.Sprintf("%+v", e.err)
}

// Unwrap returns the original underlying error.
func (e *SystemError) Unwrap() error {
	return e.err
}

// Format implements the fmt.Formatter interface.
// Implementation inspired by:
// https://github.com/pkg/errors/blob/5dd12d0cfe7f152f80558d591504ce685299311e/errors.go#L165
func (e *SystemError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", e.Cause())
			return
		}
		fallthrough
	case 's':
		// nolint: errcheck
		io.WriteString(s, e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}

// System creates a new SystemError fault whilst preserving the stack trace.
func System(pkg string, function string, msg string) *SystemError {
	m := fmt.Sprintf("%s.%s: %s", pkg, function, msg)
	return &SystemError{
		err:  pkgerrors.New(m),
		msgs: []string{m},
	}
}

// Systemf creates a new SystemError fault whilst preserving the stack trace.
func Systemf(pkg string, function string, format string, a ...interface{}) *SystemError {
	msg := fmt.Sprintf(format, a...)
	return System(pkg, function, msg)
}

// SystemWrap creates a new SystemError fault, wrapping an
// existing error and preserving the entire stack trace.
func SystemWrap(err error, pkg string, function string, msg string) *SystemError {
	var wrappedErr error
	var msgs []string

	// Purposefully using a type assertion instead of checking against all underlying errors
	// using the errors.As function so no information is lost the wrapping.
	// nolint: errorlint
	if sysErr, ok := err.(*SystemError); ok {
		pad := padding
		sb := strings.Builder{}
		for i := len(sysErr.msgs) - 1; i >= 0; i-- {
			m := sysErr.msgs[i]
			sb.WriteString(fmt.Sprintf("\n%s%s", pad, m))
			pad = pad + padding
		}
		wrappedErr = pkgerrors.Errorf("%s.%s: %s%s", pkg, function, msg, sb.String())
		msgs = sysErr.msgs
	} else {
		wrappedErr = pkgerrors.Errorf("%s.%s: %s\n%s%v", pkg, function, msg, padding, err)
		msgs = []string{err.Error()}
	}

	return &SystemError{
		err:  wrappedErr,
		msgs: append(msgs, fmt.Sprintf("%s.%s: %s", pkg, function, msg)),
	}
}

// SystemWrapf creates a new SystemError fault, wrapping an
// existing error and preserving the entire stack trace.
func SystemWrapf(
	err error,
	pkg string,
	function string,
	format string,
	a ...interface{}) *SystemError {
	return SystemWrap(err, pkg, function, fmt.Sprintf(format, a...))
}
