package fault

import (
	"errors"
	"fmt"
	"testing"
)

const (
	expectedFormat = "\n\nexpected:\n%s\n\nactual:\n%s\n\n"
)

// ------
// User Error Tests
// ------

func Test_Error_WithSingleUserError(t *testing.T) {
	code := "missing_first_name"
	msg := "Please enter your first name."
	f := User(code, msg)

	actual := f.Error()

	expected := fmt.Sprintf("%s (%s)", msg, code)
	if actual != expected {
		t.Errorf(expectedFormat, expected, actual)
	}
}

func Test_Error_WithMultipleUserErrors(t *testing.T) {
	code1 := "b"
	msg1 := "bbb"
	f := User(code1, msg1)

	code2 := "a"
	msg2 := "aaa"
	f.Add(code2, msg2)

	actual := f.Error()

	expected := fmt.Sprintf("- %s (%s)\n- %s (%s)", msg1, code1, msg2, code2)
	if actual != expected {
		t.Errorf(expectedFormat, expected, actual)
	}
}

func Test_FriendlyError_WithSingleUserError(t *testing.T) {
	code := "missing_first_name"
	msg := "Please enter your first name."
	f := User(code, msg)

	actual := f.FriendlyError()

	expected := msg
	if actual != expected {
		t.Errorf(expectedFormat, expected, actual)
	}
}

func Test_FriendlyError_WithMultipleUserErrors(t *testing.T) {
	code1 := "b"
	msg1 := "bbb"
	f := User(code1, msg1)

	code2 := "a"
	msg2 := "aaa"

	f.Add(code2, msg2)

	actual := f.FriendlyError()

	expected := fmt.Sprintf("- %s\n- %s", msg1, msg2)
	if actual != expected {
		t.Errorf(expectedFormat, expected, actual)
	}
}

func Test_Errors_WithSingleUserError(t *testing.T) {
	code := "missing_first_name"
	msg := "Please enter your first name."
	f := User(code, msg)

	actual := f.Errors()

	if len(actual) != 1 {
		t.Error("Errors() was expected to return only one key value pair.")
	}
	if actual[code] != msg {
		t.Errorf(expectedFormat, msg, actual[code])
	}
}

func Test_Errors_WithMultipleUserErrors(t *testing.T) {
	code1 := "b"
	msg1 := "bbb"
	f := User(code1, msg1)

	code2 := "a"
	msg2 := "aaa"
	f.Add(code2, msg2)

	actual := f.Errors()

	if len(actual) != 2 {
		t.Error("Errors() was expected to return two key value pairs.")
	}
	if actual[code1] != msg1 {
		t.Errorf(expectedFormat, msg1, actual[code1])
	}
	if actual[code2] != msg2 {
		t.Errorf(expectedFormat, msg2, actual[code2])
	}
}

func Test_ErrorMessages_WithSingleUserError(t *testing.T) {
	code := "missing_first_name"
	msg := "Please enter your first name."
	f := User(code, msg)

	actual := f.ErrorMessages()

	if len(actual) != 1 {
		t.Error("ErrorMessages() was expected to return only one message.")
	}
	if actual[0] != msg {
		t.Errorf(expectedFormat, msg, actual[0])
	}
}

func Test_ErrorMessages_WithMultipleUserErrors(t *testing.T) {
	code1 := "b"
	msg1 := "bbb"
	f := User(code1, msg1)

	code2 := "a"
	msg2 := "aaa"
	f.Add(code2, msg2)

	actual := f.ErrorMessages()

	if len(actual) != 2 {
		t.Error("ErrorMessages() was expected to return two messages.")
	}
	if actual[0] != msg1 {
		t.Errorf(expectedFormat, msg1, actual[0])
	}
	if actual[1] != msg2 {
		t.Errorf(expectedFormat, msg2, actual[1])
	}
}

// ------
// System Error Tests
// ------

func Test_Error_WithSingleSystemError(t *testing.T) {
	f := System("a", "b", "c")

	actual := f.Error()

	expected := "a.b: c"
	if actual != expected {
		t.Errorf(expectedFormat, expected, actual)
	}
}

func Test_String_WithSingleSystemError(t *testing.T) {
	f := System("a", "b", "c")

	actual := f.String()

	expected := "a.b: c"
	if actual != expected {
		t.Errorf(expectedFormat, expected, actual)
	}
}

func Test_Error_WithLayersOfSystemErrorsAndOneNonSystemError(t *testing.T) {
	f1 := errors.New("foo bar")
	f2 := SystemWrap(f1, "d", "e", "f")
	f3 := SystemWrap(f2, "g", "h", "i")

	actual := f3.Error()

	expected := "g.h: i\n   d.e: f\n      foo bar"
	if actual != expected {
		t.Errorf(expectedFormat, expected, actual)
	}
}

func Test_String_WithLayersOfSystemErrorsAndOneNonSystemError(t *testing.T) {
	f1 := errors.New("foo bar")
	f2 := SystemWrap(f1, "d", "e", "f")
	f3 := SystemWrap(f2, "g", "h", "i")

	actual := f3.String()

	expected := "g.h: i\n   d.e: f\n      foo bar"
	if actual != expected {
		t.Errorf(expectedFormat, expected, actual)
	}
}

func Test_Error_WithLayersOfSystemErrors(t *testing.T) {
	f1 := System("a", "b", "c")
	f2 := SystemWrap(f1, "d", "e", "f")
	f3 := SystemWrap(f2, "g", "h", "i")

	actual := f3.Error()

	expected := "g.h: i\n   d.e: f\n      a.b: c"
	if actual != expected {
		t.Errorf(expectedFormat, expected, actual)
	}
}

func Test_FormatWithoutPlus_WithLayersOfSystemErrors_ReturnsSameAsError(t *testing.T) {
	f1 := System("a", "b", "c")
	f2 := SystemWrap(f1, "d", "e", "f")
	f3 := SystemWrap(f2, "g", "h", "i")

	expected := f3.Error()
	notExpected := f3.StackTrace()

	actual := fmt.Sprintf("%v", f3)
	if actual != expected {
		t.Errorf(expectedFormat, expected, actual)
	}
	if actual == notExpected {
		t.Error("The 'v' formatter should not include a stack trace.")
	}
}

func Test_FormatWithPlus_WithLayersOfSystemErrors_ReturnsSameAsStackTrace(t *testing.T) {
	f1 := System("a", "b", "c")
	f2 := SystemWrap(f1, "d", "e", "f")
	f3 := SystemWrap(f2, "g", "h", "i")

	expected := f3.StackTrace()
	notExpected := f3.Error()

	actual := fmt.Sprintf("%+v", f3)
	if actual != expected {
		t.Errorf(expectedFormat, expected, actual)
	}
	if actual == notExpected {
		t.Error("The '+v' formatter should include a stack trace.")
	}
}
