package bloom

import (
	"testing"
)

func TestBloomFilter(t *testing.T) {
	bf := NewBloom()

	bf.Add("golang")
	bf.Add("bloom")
	bf.Add("filter")

	if !bf.Contains("golang") {
		t.Errorf("expected 'golang' in bloom filter")
	}
	if !bf.Contains("bloom") {
		t.Errorf("expected 'bloom' in bloom filter")
	}
	if bf.Contains("bitcoin") {
		t.Log("false positive")
	}
}

func TestFalsePositiveRate(t *testing.T) {
	bf := NewBloom()

	// Add 1000 items
	for i := range 1000 {
		bf.Add(string(rune(i)))
	}

	// Check another 1000 different items
	falsePositives := 0
	for i := 1000; i < 2000; i++ {
		if bf.Contains(string(rune(i))) {
			falsePositives++
		}
	}

	t.Logf("false positive number: %d", falsePositives)

	rate := float64(falsePositives) / 1000.0
	t.Logf("False positive rate observed: %f", rate)

	if rate > 0.01 {
		t.Errorf("False positive rate too high: %f", rate)
	}
}
