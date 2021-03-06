package email

import (
	"testing"
)

func Test_New_Unhappy(t *testing.T) {
	cases := []string{
		"",
		"a",
		"a.b",
		"foo@bar",
		"a@b.c",
		"foo.com@bar",
	}

	for _, value := range cases {
		_, err := New(value)

		if err == nil {
			t.Errorf("The value '%s' was expected to fail email validation.", value)
		}
	}
}

func Test_New_Happy(t *testing.T) {
	cases := []string{
		"foo@bar.com",
		"a@b.io",
	}

	for _, value := range cases {
		_, err := New(value)

		if err != nil {
			t.Errorf("The value '%s' failed email validation.", value)
		}
	}
}

func Test_Domain_Normalised(t *testing.T) {
	cases := []string{
		"foo@bar.com",
		"foo@bar.com ",
		"  foo@bar.com",
		"  foo@bar.com    ",
		"FOO@bar.com",
		"FOO@BAR.COM",
		"fOo@bAr.coM",
	}

	for _, value := range cases {
		addr, err := New(value)

		if err != nil {
			t.Errorf("The value '%s' failed email validation.", value)
		}

		if addr.Domain() != "bar.com" {
			t.Errorf("Domain was expected to return 'bar.com' for value '%s', but returned '%s'.", value, addr.Domain())
		}

		if addr.Normalised() != "foo@bar.com" {
			t.Errorf("Normalised was expected to return 'foo@bar.com' for value '%s', but returned '%s'.", value, addr.Normalised())
		}
	}
}
