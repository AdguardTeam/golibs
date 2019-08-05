package utils

import "testing"

func TestHostValidation(t *testing.T) {
	validHosts := []string{
		"0example.com", "example.com", "ex.example.com", "ex-1ample.com.ru",
		"xn---asdasd.com", "local", "aa.ru", "a.ru", "00.11.22.33",
	}
	invalidHosts := []string{
		"-a.c", "a-.c", "a.-c", "a.c-",
		"host-", "h@st", "*.com", "ex_ample.com", "!asd.ru", "google..com",
		".google.com", "google.com.",
	}

	for _, h := range validHosts {
		if err := IsValidHostname(h); err != nil {
			t.Fatalf("host %s is valid, but IsValidHostname returns error: %s", h, err)
		}
	}

	for _, h := range invalidHosts {
		if err := IsValidHostname(h); err == nil {
			t.Fatalf("host %s is invalid, but IsValidHostname returns nil", h)
		}
	}
}
