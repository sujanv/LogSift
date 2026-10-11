package geodb

import "net"

type CIDRTrie struct{}

func NewCIDRTrie() *CIDRTrie { return &CIDRTrie{} }
func (c *CIDRTrie) Insert(cidr, code, name string) {}
func (c *CIDRTrie) Search(ip net.IP) (string, string) { return "", "" }
