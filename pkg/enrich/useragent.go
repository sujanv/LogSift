package enrich

import (
	"strings"
)

type UserAgentInfo struct {
	Browser string `json:"browser"`
	OS      string `json:"os"`
	Device  string `json:"device"`
	IsBot   bool   `json:"is_bot"`
}

func ParseUserAgent(ua string) UserAgentInfo {
	lower := strings.ToLower(ua)
	info := UserAgentInfo{Browser: "Other", OS: "Other", Device: "Desktop"}

	if strings.Contains(lower, "bot") || strings.Contains(lower, "spider") ||
		strings.Contains(lower, "crawl") || strings.Contains(lower, "slurp") {
		info.IsBot = true
	}

	switch {
	case strings.Contains(lower, "edg"):
		info.Browser = "Edge"
	case strings.Contains(lower, "chrome") && !strings.Contains(lower, "chromium"):
		info.Browser = "Chrome"
	case strings.Contains(lower, "firefox"):
		info.Browser = "Firefox"
	case strings.Contains(lower, "safari") && !strings.Contains(lower, "chrome"):
		info.Browser = "Safari"
	case strings.Contains(lower, "curl"):
		info.Browser = "cURL"
		info.IsBot = true
	}

	switch {
	case strings.Contains(lower, "iphone") || strings.Contains(lower, "ipad"):
		info.OS = "iOS"
		info.Device = "Mobile"
	case strings.Contains(lower, "android"):
		info.OS = "Android"
		info.Device = "Mobile"
	case strings.Contains(lower, "windows"):
		info.OS = "Windows"
	case strings.Contains(lower, "macintosh") || strings.Contains(lower, "mac os x"):
		info.OS = "macOS"
	case strings.Contains(lower, "linux"):
		info.OS = "Linux"
	}

	return info
}
