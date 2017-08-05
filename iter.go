package shufn

import (
	"math/rand"
)

// New creates a non-thread-safe iterator over the numeric range.
func New(mult, mod, min, max, start uint64) *Iter {
	if start == 0 {
		start = rand.Uint64()
	}
	start = start % (max - min)
	if start == 0 {
		start = 1
	}

	return &Iter{
		Mult:  mult,
		Mod:   mod,
		Start: start,
		Min:   min,
		Max:   max,
	}
}

type Iter struct {
	Mult  uint64
	Mod   uint64
	Start uint64
	i     uint64
	I     uint64
	Max   uint64
	Min   uint64
}

// Next returns whether there are more numbers in the sequence; and indicates
// the next is available on i.I. (Not thread safe!)
func (i *Iter) Next() (more bool) {
	i.I, more = i.NextI()
	return
}

func (i *Iter) NextI() (I uint64, more bool) {
	I, more = i.next()
	for more && I > i.Max {
		I, more = i.next()
	}
	return
}

func (i *Iter) next() (I uint64, more bool) {
	if i.i == 0 {
		i.i = i.Start
		more = true
	} else {
		i.i = (i.i * i.Mult) % i.Mod
		more = i.i != i.Start
	}

	I = i.i + i.Min - 1
	return
}
