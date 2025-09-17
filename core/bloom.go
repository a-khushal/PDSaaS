package main

import (
	"fmt"
	"math"
)

const n = 1_000_000 // number of items in the filter
const p = 0.001     // false positive porbability that we allow

func main() {
	m := int(math.Ceil(-float64(n) * math.Log(p) / math.Pow(math.Log(2), 2))) // the size of the bit array
	k := int(math.Ceil((float64(m) / float64(n)) * math.Log(2)))              // Optimal number of hash functions

	filter := make([]int, m)

	fmt.Println(filter)
}
