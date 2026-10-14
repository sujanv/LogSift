package crypto

import (
	"testing"
)

func TestTamperChain(t *testing.T) {
	tc := NewTamperChain("super-secret-audit-key")

	l1 := []byte("2026-10-14 10:00:00 user logged in")
	l2 := []byte("2026-10-14 10:01:00 transaction 500 dollars")
	l3 := []byte("2026-10-14 10:02:00 user logged out")

	tc.Append(l1)
	tc.Append(l2)
	tc.Append(l3)

	valid, count, err := tc.Verify([][]byte{l1, l2, l3})
	if !valid || err != nil || count != 3 {
		t.Fatalf("chain must be valid: %v", err)
	}

	tamperedL2 := []byte("2026-10-14 10:01:00 transaction 5000000 dollars")
	valid, badIndex, err := tc.Verify([][]byte{l1, tamperedL2, l3})
	if valid || badIndex != 2 {
		t.Errorf("tamper must be detected at block 2, got index %d, err %v", badIndex, err)
	}
}
