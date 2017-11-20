package shufn_test

import (
	. "github.com/icio/shufn"

	"fmt"
	"math/rand"
	"time"
)

// Here we produce all of the numbers 1-10 in a slightly random order.
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

// By specifying Mult, Mod, and Loop we always get the same sequence!
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

// To resume a consistent sequence, just set I before consuming.
func ExampleIter_resume() {
	rand.Seed(time.Now().Unix())

	iter := MustIter(&Seq{
		Min: 100, Max: 110,
		Mult: 7, Mod: 11, Loop: 106,
	})
	iter.I = 109 // Resume, starting from 109.

	for iter.Next() {
		fmt.Print(iter.I, " ")
	}
	// Output: 109 103 105 108 107 100
}
