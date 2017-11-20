package shufn

import (
	"fmt"
	"testing"
)

type primesPastFunc func(uint64) []uint64

func BenchmarkPrimesPast(b *testing.B) {
	impl := map[string]primesPastFunc{
		"append":      primesPast_append,
		"overalloc":   primesPast_overalloc,
		"alloczeroes": primesPast_alloczeroes,
	}
	ns := []uint64{1, 10, 100, 1000, 10000, 100000, 1000000}

	subTest := func(fn primesPastFunc, n uint64) func(*testing.B) {
		return func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				fn(n)
			}
		}
	}

	for _, n := range ns {
		for name, fn := range impl {
			b.Run(
				fmt.Sprintf("%s_%d", name, n),
				subTest(fn, n),
			)
		}
	}
}
