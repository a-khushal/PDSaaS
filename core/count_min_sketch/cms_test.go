package countminsketch

import (
	"fmt"
	"math"
	"testing"
)

func TestBasicCountsWithBounds(t *testing.T) {
	epsilon := 0.001
	delta := 1e-9
	cms := NewCountMinSketch(epsilon, delta)

	trueApple := uint32(5)
	trueBanana := uint32(3)
	trueCherry := uint32(0)

	for range trueApple {
		cms.Add("apple")
	}
	for range trueBanana {
		cms.Add("banana")
	}

	apple := cms.Count("apple")
	banana := cms.Count("banana")
	cherry := cms.Count("cherry")

	if apple < trueApple || banana < trueBanana || cherry < trueCherry {
		t.Fatalf("undercount detected: got apple=%d banana=%d cherry=%d", apple, banana, cherry)
	}

	total := float64(trueApple + trueBanana + trueCherry)
	maxError := uint32(math.Ceil(epsilon * total))
	if apple > trueApple+maxError {
		t.Fatalf("apple overcount too large: got=%d true=%d maxError=%d", apple, trueApple, maxError)
	}
	if banana > trueBanana+maxError {
		t.Fatalf("banana overcount too large: got=%d true=%d maxError=%d", banana, trueBanana, maxError)
	}
}

func TestNoUndercount(t *testing.T) {
	cms := NewCountMinSketch(0.001, 1e-9)

	for i := range 1000 {
		cms.Add("x")
		if i%2 == 0 {
			cms.Add("y")
		}
	}

	if got := cms.Count("x"); got < 1000 {
		t.Fatalf("undercount for x: got=%d want>=1000", got)
	}
	if got := cms.Count("y"); got < 500 {
		t.Fatalf("undercount for y: got=%d want>=500", got)
	}
}

func TestHeavyHitterDominates(t *testing.T) {
	cms := NewCountMinSketch(0.001, 1e-9)

	for range 1000 {
		cms.Add("heavy")
	}

	for i := range 100 {
		key := fmt.Sprintf("light_%d", i)
		for range 5 {
			cms.Add(key)
		}
	}

	heavy := cms.Count("heavy")
	var maxLight uint32
	for i := range 100 {
		key := fmt.Sprintf("light_%d", i)
		if c := cms.Count(key); c > maxLight {
			maxLight = c
		}
	}

	if heavy <= maxLight {
		t.Fatalf("heavy hitter not dominating: heavy=%d maxLight=%d", heavy, maxLight)
	}
}

func TestAbsentIsZero(t *testing.T) {
	cms := NewCountMinSketch(0.001, 1e-9)
	if got := cms.Count("missing"); got != 0 {
		t.Fatalf("absent key should be 0, got=%d", got)
	}
	cms.Add("other")
	if got := cms.Count("missing"); got != 0 {
		t.Fatalf("absent key should remain 0 after other inserts, got=%d", got)
	}
}
