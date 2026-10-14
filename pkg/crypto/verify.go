package crypto

func (tc *TamperChain) Blocks() []LogBlock {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	res := make([]LogBlock, len(tc.blocks))
	copy(res, tc.blocks)
	return res
}
