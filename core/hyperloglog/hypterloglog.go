package hyperLogLog

import (
	"math"

	"github.com/spaolacci/murmur3"
)

type HyperLogLog struct {
	p      uint8   // precision (14-16)
	m      uint32  // number of registers (2^p)
	alphaM float64 // correction constant
	M      []uint8 // registers
}

func NewHyperLogLog(p uint8) *HyperLogLog {
	m := uint32(1) << p
	var alphaM float64

	switch m {
	case 16:
		alphaM = 0.673
	case 32:
		alphaM = 0.697
	case 64:
		alphaM = 0.709
	default:
		alphaM = 0.7213 / (1 + 1.079/float64(m))
	}

	return &HyperLogLog{
		p:      p,
		m:      m,
		alphaM: alphaM,
		M:      make([]uint8, m),
	}
}

// rho counts the position of the first 1-bit
// it is the index of the leftmost 1-bit in w
func rho(w uint32, max uint8) uint8 {
	var r uint8 = 1
	for (w&(1<<31)) == 0 && r <= max {
		r++
		w <<= 1
	}
	return r
}

// Count estimates the cardinality
func (hll *HyperLogLog) Add(item string) {
	x := murmur3.Sum32([]byte(item)) // output: [0 .. 2³² - 1]
	idx := x >> (32 - hll.p)         // first p bits, we do 32 cuz of above output range
	w := x << hll.p                  // remaining bits
	r := rho(w, uint8(32-hll.p))     // streak of zeros = index of the leftmost 1-bit in w

	/*
	* Shorter zero streaks (small r) are common -> not very informative.
	* Longer zero streaks (big r) are rare -> strong evidence of high cardinality.
	* HLL relies on the maximum leading-zero count per register as the key statistic.
	 */
	if r > hll.M[idx] {
		hll.M[idx] = r
	}
}

/*
	Raw HLL estimate:

* E = αm * m^2 * (1 / ∑(2^-M[j]))

* sum = ∑(2^-M[j])         -> denominator
* Z   = 1 / sum            -> reciprocal denominator
* E   = αm * m^2 * Z       -> scaled harmonic mean of 2^M[j](see the denominator above) -> cardinolity
*/
func (hll *HyperLogLog) Count() int {
	var sum float64
	var zeros uint32
	for _, v := range hll.M {
		sum += 1.0 / float64(uint64(1)<<v)
		if v == 0 {
			zeros++
		}
	}
	Z := 1.0 / sum
	E := hll.alphaM * float64(hll.m*hll.m) * Z

	if zeros > 0 {
		smallEstimate := float64(hll.m) * math.Log(float64(hll.m)/float64(zeros))
		if E <= (5.0/2.0)*float64(hll.m) {
			return int(smallEstimate)
		}
	}

	return int(E)
}
