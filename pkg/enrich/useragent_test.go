package enrich

import (
	"testing"
)

func TestParseUserAgent(t *testing.T) {
	ua1 := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36"
	info1 := ParseUserAgent(ua1)
	if info1.Browser != "Chrome" || info1.OS != "macOS" || info1.IsBot {
		t.Errorf("mismatch on chrome mac: %+v", info1)
	}

	ua2 := "Googlebot/2.1 (+http://www.google.com/bot.html)"
	info2 := ParseUserAgent(ua2)
	if !info2.IsBot {
		t.Errorf("Googlebot must be detected as bot")
	}
}
