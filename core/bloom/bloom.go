package bloom

import (
	"math"

	city "github.com/go-faster/city"
	"github.com/spaolacci/murmur3"
)

const n = 1_000_000 // expected items
const p = 0.001     // false positive probability

type Bloom struct {
	m      int
	k      int
	bitset []bool
}

// creates a Bloom filter.
func NewBloom() *Bloom {
	m := int(math.Ceil(-float64(n) * math.Log(p) / math.Pow(math.Log(2), 2)))
	k := int(math.Ceil((float64(m) / float64(n)) * math.Log(2)))
	if k < 1 {
		k = 1
	}
	return &Bloom{
		m:      m,
		k:      k,
		bitset: make([]bool, m),
	}
}

func (bf *Bloom) hashes(data []byte) []int {
	if bf.m == 0 || bf.k == 0 {
		return nil
	}

	hashes := make([]int, 0, bf.k)

	// base hashes both are in range [0, m-1]
	h1 := int(murmur3.Sum32(data)) % bf.m
	h2 := int(city.Hash64(data) % uint64(bf.m))

	// ensure non-negative (Sum32 -> uint32 cast to int is non-negative, but safe-guard anyway)
	if h1 < 0 {
		h1 += bf.m
	}
	if h2 < 0 {
		h2 += bf.m
	}
	// If h2 is 0, using (h1 + i*h2) will produce duplicates; make sure h2 != 0
	if h2 == 0 {
		h2 = 1
	}
	for i := 0; i < bf.k; i++ {
		idx := (h1 + i*h2) % bf.m
		hashes = append(hashes, idx)
	}
	return hashes
}

func (bf *Bloom) Add(item string) {
	for _, idx := range bf.hashes([]byte(item)) {
		bf.bitset[idx] = true
	}
}

func (bf *Bloom) Contains(item string) bool {
	for _, idx := range bf.hashes([]byte(item)) {
		if !bf.bitset[idx] {
			return false
		}
	}
	return true
}
