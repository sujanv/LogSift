package geodb

import (
	"net"
	"strings"
)

type GeoInfo struct {
	CountryCode string `json:"country_code"`
	CountryName string `json:"country_name"`
	IsPrivate   bool   `json:"is_private"`
	ASN         string `json:"asn,omitempty"`
}

type Resolver struct {
	trie *CIDRTrie
}

func NewResolver() *Resolver {
	t := NewCIDRTrie()
	t.Insert("10.0.0.0/8", "PRIVATE", "Private Network")
	t.Insert("172.16.0.0/12", "PRIVATE", "Private Network")
	t.Insert("192.168.0.0/16", "PRIVATE", "Private Network")
	t.Insert("127.0.0.0/8", "LOOPBACK", "Localhost Loopback")
	t.Insert("169.254.0.0/16", "LINKLOCAL", "Link-Local")

	t.Insert("8.8.8.0/24", "US", "United States (Google DNS)")
	t.Insert("1.1.1.0/24", "AU", "Australia (Cloudflare)")
	t.Insert("13.0.0.0/8", "US", "United States (AWS)")
	t.Insert("185.0.0.0/8", "EU", "European Union")
	t.Insert("203.0.0.0/8", "AP", "Asia Pacific")

	return &Resolver{trie: t}
}

func (r *Resolver) Lookup(ipStr string) GeoInfo {
	ipStr = strings.TrimSpace(ipStr)
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return GeoInfo{CountryCode: "UNKNOWN", CountryName: "Unknown IP"}
	}

	if ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast() {
		return GeoInfo{
			CountryCode: "LAN",
			CountryName: "Private/Local Network",
			IsPrivate:   true,
		}
	}

	code, name := r.trie.Search(ip)
	if code != "" {
		return GeoInfo{
			CountryCode: code,
			CountryName: name,
			IsPrivate:   code == "PRIVATE" || code == "LOOPBACK",
		}
	}

	return GeoInfo{
		CountryCode: "GLOBAL",
		CountryName: "Public Internet",
		IsPrivate:   false,
	}
}
