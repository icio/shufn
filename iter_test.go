package shufn_test

import (
	. "github.com/icio/shufn"

	"fmt"
	"math/rand"
	"strings"
	"testing"
)

func TestIter(t *testing.T) {
	seed := int64(100)
	tests := []struct {
		seq *Seq
		out []uint64
		err error
	}{
		{&Seq{Min: 6, Max: 5}, nil, ErrInvalidBounds},
		{&Seq{Min: 5, Max: 6, Loop: 4}, nil, ErrInvalidLoop},
		{&Seq{Min: 5, Max: 5}, []uint64{5}, nil},                                    // Min==Max.
		{&Seq{Min: 5, Max: 6, Loop: 5}, []uint64{5, 6}, nil},                        // Max-Min==1, Loop==Min.
		{&Seq{Min: 5, Max: 6, Loop: 6}, []uint64{6, 5}, nil},                        // Max-Min==1, Loop==Max.
		{&Seq{Min: 4, Max: 7}, []uint64{4, 6, 7, 5}, nil},                           // Simple, random iterator.
		{&Seq{Min: 4, Max: 7, Loop: 5}, []uint64{5, 4, 6, 7}, nil},                  // Random iterator with fixed starting point.
		{&Seq{Min: 4, Max: 7, Mult: 3, Mod: 5, Loop: 5}, []uint64{5, 4, 6, 7}, nil}, // Consistent random iterator.
	}

	for n, test := range tests {
		t.Run(testName(n, test.seq), func(t *testing.T) {
			// Actual.
			rand.Seed(seed)
			iter, err := NewIter(test.seq)

			// Check the error.
			if exp, act := test.err, err; exp != act {
				t.Errorf("Expected error %s but got %s.", errQuote(exp), errQuote(act))
				return
			}
			if err != nil {
				return
			}

			// Check the iter.
			t.Logf("iter: %#v, iter.Seq: %#v", iter, iter.Seq)
			// fmt.Printf("iter: %#v, iter.Seq: %#v\n", iter, iter.Seq)
			if exp, act := fmt.Sprintf("%v", test.out), fmt.Sprintf("%v", Consume(iter)); exp != act {
				t.Errorf("Expected sequence %s but got %s.", exp, act)
			}
		})
	}
}

func errQuote(err error) string {
	if err == nil {
		return fmt.Sprintf("%v", err)
	}
	return fmt.Sprintf("%q", err)
}

func testName(n int, s *Seq) string {
	return fmt.Sprintf("%d:%s", n, strings.Replace(fmt.Sprintf("%+v", s), " ", ",", -1))
}
