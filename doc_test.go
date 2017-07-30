package shufn_test

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/icio/shufn"
)

func ExampleNew() {
	// it produces values in the range 0-100 inclusive. The start parameter is
	// passed through to New verbatim. When New receives start=0 it replaces it
	// with a random value.
	rand.Seed(0)
	it := shufn.New(shufn.Calc(1, 3, 0))
	for {
		i, more := it.Next()
		if !more {
			break
		}

		fmt.Println(i)
	}
	// Output: 1
	// 3
	// 2
}
