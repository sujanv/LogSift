package geodb

import (
	"testing"
)

func TestResolver(t *testing.T) {
	res := NewResolver()

	info1 := res.Lookup("192.168.1.100")
	if !info1.IsPrivate {
		t.Errorf("192.168.1.100 must be marked private")
	}

	info2 := res.Lookup("127.0.0.1")
	if !info2.IsPrivate {
		t.Errorf("127.0.0.1 must be marked private")
	}

	info3 := res.Lookup("8.8.8.8")
	if info3.CountryCode != "US" {
		t.Errorf("expected US for 8.8.8.8, got %s", info3.CountryCode)
	}

	info4 := res.Lookup("1.1.1.1")
	if info4.CountryCode != "AU" {
		t.Errorf("expected AU for 1.1.1.1, got %s", info4.CountryCode)
	}
}
