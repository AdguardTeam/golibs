package utils

import "testing"

func TestHostValidation(t *testing.T) {
	validHosts := []string{"example.com", "ex.example.com", "ex-1ample.com.ru", "xn---asdasd.com", "local", "aa.ru"}
	invalidHosts := []string{"-a.c", "host-", "h@st", "*.com", "ex_ample.com", "!asd.ru", "google..com"}

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
