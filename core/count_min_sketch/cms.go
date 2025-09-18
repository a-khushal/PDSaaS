package countminsketch

import (
	"encoding/binary"
	"math"
	"math/rand/v2"

	"github.com/spaolacci/murmur3"
)

type CountMinSketch struct {
	width  uint32     // number of columns
	depth  uint32     // number of rows
	counts [][]uint32 // 2D array of counters
	seeds  []uint32   // seeds for each hash function
}

func NewCountMinSketch(epsilon, delta float64) *CountMinSketch {
	w := uint32(math.Ceil(math.E / epsilon))
	d := uint32(math.Ceil(math.Log(1 / delta)))

	counts := make([][]uint32, d)
	for i := range counts {
		counts[i] = make([]uint32, w)
	}

	seeds := make([]uint32, d)
	for i := range seeds {
		seeds[i] = rand.Uint32()
	}

	return &CountMinSketch{
		width:  w,
		depth:  d,
		counts: counts,
		seeds:  seeds,
	}
}

// Add increments the count for an item
func (cms *CountMinSketch) Add(item string) {
	for i := uint32(0); i < cms.depth; i++ {
		idx := cms.hash(item, cms.seeds[i]) % cms.width
		cms.counts[i][idx]++
	}
}

// Count estimates the frequency of an item
func (cms *CountMinSketch) Count(item string) uint32 {
	min := uint32(math.MaxUint32)
	for i := uint32(0); i < cms.depth; i++ {
		idx := cms.hash(item, cms.seeds[i]) % cms.width
		if cms.counts[i][idx] < min {
			min = cms.counts[i][idx]
		}
	}
	return min
}

func (cms *CountMinSketch) hash(item string, seed uint32) uint32 {
	data := []byte(item)

	h := murmur3.New32()
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, seed)
	h.Write(b)
	h.Write(data)

	return h.Sum32()
}
