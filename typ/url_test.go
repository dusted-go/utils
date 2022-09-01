package typ

import "testing"

func Test_New_WithoutHttpOrHttps(t *testing.T) {
	v := "www.example.org"

	url := NewURL(v)

	if url.String() != "https://www.example.org" {
		t.Errorf("New created wrong URL format: %s", url.String())
	}
}

func Test_New_WithHttp(t *testing.T) {
	v := "http://www.example.org"

	url := NewURL(v)

	if url.String() != "http://www.example.org" {
		t.Errorf("New created wrong URL format: %s", url.String())
	}
}
