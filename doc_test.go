package shufn_test

import (
	. "github.com/icio/shufn"

	"fmt"
	"math/rand"
	"time"
)

// iter enumerates the integers in the range [1, 10] pseudo-randomly
// with no guarantees.
func ExampleIter() {
	rand.Seed(0) // Ensure reproducible results in examples. Set to random.

	iter := MustIter(&Seq{
		Min: 1,
		Max: 10,
	})
	for iter.Next() {
		fmt.Print(iter.I, " ")
	}
	// Output: 3 10 4 6 9 8 1 7 5 2
}

// iter enumerates the integers in the range [1, 10], always in the same
// sequence. Mult, Mod and Loop are the only factors randomly set if zero.
func ExampleIter_consistent() {
	rand.Seed(time.Now().Unix())

	iter := MustIter(&Seq{
		Min: 1, Max: 10,
		Mult: 6, Mod: 11, Loop: 6,
	})
	for iter.Next() {
		fmt.Print(iter.I, " ")
	}
	// Output: 6 3 7 9 10 5 8 4 2 1
}

func ExampleIter_resume() {
	rand.Seed(time.Now().Unix())

	iter := MustIter(&Seq{
		Min: 100, Max: 110,
		Mult: 7, Mod: 11, Loop: 106,
	})

	// Resume from 10. We could alternatively have set Start on the Seq passed
	// to MustIter, prior to its invocation.
	iter.I = 109

	for iter.Next() {
		fmt.Print(iter.I, " ")
	}
	// Output: 109 103 105 108 107 100
}
