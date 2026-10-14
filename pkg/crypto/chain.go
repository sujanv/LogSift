package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
)

type LogBlock struct {
	Index     int64  `json:"index"`
	PrevHash  string `json:"prev_hash"`
	BlockHash string `json:"block_hash"`
	Signature string `json:"signature,omitempty"`
}

type TamperChain struct {
	mu       sync.Mutex
	secret   []byte
	prevHash string
	blocks   []LogBlock
}

func NewTamperChain(secretKey string) *TamperChain {
	return &TamperChain{
		secret:   []byte(secretKey),
		prevHash: "0000000000000000000000000000000000000000000000000000000000000000",
		blocks:   make([]LogBlock, 0, 100),
	}
}

func (tc *TamperChain) Append(data []byte) LogBlock {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	idx := int64(len(tc.blocks) + 1)
	hasher := sha256.New()
	hasher.Write([]byte(tc.prevHash))
	hasher.Write(data)
	blockHash := hex.EncodeToString(hasher.Sum(nil))

	var sig string
	if len(tc.secret) > 0 {
		mac := hmac.New(sha256.New, tc.secret)
		mac.Write([]byte(blockHash))
		sig = hex.EncodeToString(mac.Sum(nil))
	}

	block := LogBlock{
		Index:     idx,
		PrevHash:  tc.prevHash,
		BlockHash: blockHash,
		Signature: sig,
	}

	tc.blocks = append(tc.blocks, block)
	tc.prevHash = blockHash
	return block
}

func (tc *TamperChain) Verify(entries [][]byte) (bool, int64, error) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	if len(entries) != len(tc.blocks) {
		return false, 0, fmt.Errorf("length mismatch: %d entries vs %d blocks", len(entries), len(tc.blocks))
	}

	prev := "0000000000000000000000000000000000000000000000000000000000000000"
	for i, b := range tc.blocks {
		if b.PrevHash != prev {
			return false, b.Index, fmt.Errorf("broken link at block %d", b.Index)
		}
		hasher := sha256.New()
		hasher.Write([]byte(prev))
		hasher.Write(entries[i])
		expectedHash := hex.EncodeToString(hasher.Sum(nil))
		if b.BlockHash != expectedHash {
			return false, b.Index, fmt.Errorf("tamper detected at block %d: hash mismatch", b.Index)
		}
		prev = b.BlockHash
	}
	return true, int64(len(tc.blocks)), nil
}
