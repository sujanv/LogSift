package archive

import (
	"hash/fnv"
)

type ChunkBloomFilter struct {
	bits []uint64
	size uint64
}

func NewChunkBloomFilter(bitSize uint64) *ChunkBloomFilter {
	if bitSize == 0 {
		bitSize = 1024
	}
	words := (bitSize + 63) / 64
	return &ChunkBloomFilter{
		bits: make([]uint64, words),
		size: bitSize,
	}
}

func (bf *ChunkBloomFilter) Add(key string) {
	h1 := hashFnv(key, 0) % bf.size
	h2 := hashFnv(key, 1) % bf.size
	h3 := hashFnv(key, 2) % bf.size

	bf.setBit(h1)
	bf.setBit(h2)
	bf.setBit(h3)
}

func (bf *ChunkBloomFilter) Test(key string) bool {
	h1 := hashFnv(key, 0) % bf.size
	h2 := hashFnv(key, 1) % bf.size
	h3 := hashFnv(key, 2) % bf.size

	return bf.getBit(h1) && bf.getBit(h2) && bf.getBit(h3)
}

func (bf *ChunkBloomFilter) setBit(idx uint64) {
	word := idx / 64
	pos := idx % 64
	bf.bits[word] |= (1 << pos)
}

func (bf *ChunkBloomFilter) getBit(idx uint64) bool {
	word := idx / 64
	pos := idx % 64
	return (bf.bits[word] & (1 << pos)) != 0
}

func hashFnv(data string, seed uint8) uint64 {
	h := fnv.New64a()
	h.Write([]byte{seed})
	h.Write([]byte(data))
	return h.Sum64()
}
