package geodb

import "net"

type cidrEntry struct {
	network *net.IPNet
	code    string
	name    string
}

type CIDRTrie struct {
	entries []cidrEntry
}

func NewCIDRTrie() *CIDRTrie {
	return &CIDRTrie{entries: make([]cidrEntry, 0, 32)}
}

func (t *CIDRTrie) Insert(cidr, code, name string) {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return
	}
	t.entries = append(t.entries, cidrEntry{
		network: ipnet,
		code:    code,
		name:    name,
	})
}

func (t *CIDRTrie) Search(ip net.IP) (string, string) {
	for _, e := range t.entries {
		if e.network.Contains(ip) {
			return e.code, e.name
		}
	}
	return "", ""
}
