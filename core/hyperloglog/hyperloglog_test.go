package hyperLogLog

import (
	"fmt"
	"math"
	"testing"
)

func TestBasicEstimate(t *testing.T) {
	hll := NewHyperLogLog(14) // p=14 -> m = 2^14

	for i := 0; i < 10000; i++ {
		hll.Add(fmt.Sprintf("user_%d", i))
	}

	est := hll.Count()
	trueN := 10000.0 // real number of unique elements in the array
	relErr := math.Abs(trueN-float64(est)) / trueN

	// HLL's error is roughly 1.04/sqrt(m)
	maxRelErr := 0.02
	if relErr > maxRelErr {
		t.Fatalf("relative error too high: got=%f want<=%f (estimate=%d)", relErr, 0.05, est)
	}
}

func TestDuplicatesDoNotIncreaseCount(t *testing.T) {
	hll := NewHyperLogLog(14)

	for i := 0; i < 5000; i++ {
		hll.Add(fmt.Sprintf("item_%d", i))
	}
	first := hll.Count()

	for i := 0; i < 5000; i++ {
		hll.Add(fmt.Sprintf("item_%d", i))
	}
	second := hll.Count()

	if second < first {
		t.Fatalf("count decreased after duplicates: first=%d second=%d", first, second)
	}

	trueN := 5000.0
	relErr := math.Abs(trueN-float64(second)) / trueN
	if relErr > 0.02 {
		t.Fatalf("relative error too high with duplicates: got=%f want<=%f (estimate=%d)", relErr, 0.10, second)
	}
}

func TestMonotonicity(t *testing.T) {
	hll := NewHyperLogLog(14)

	for i := 0; i < 1000; i++ {
		hll.Add(fmt.Sprintf("k_%d", i))
	}
	c1 := hll.Count()

	for i := 0; i < 10000; i++ {
		hll.Add(fmt.Sprintf("k_%d", i))
	}
	c2 := hll.Count()

	if c2 < c1 {
		t.Fatalf("estimate decreased after adding more items: before=%d after=%d", c1, c2)
	}
}

func TestEmptyIsZero(t *testing.T) {
	hll := NewHyperLogLog(14)
	if got := hll.Count(); got != 0 {
		t.Fatalf("empty HLL should be 0, got=%d", got)
	}
}
