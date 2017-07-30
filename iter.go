package shufn

import (
	"math/rand"
)

// Iter defines the common interface to thread-safe and -unsafe variants
// of the iterator.
type Iter interface {
	Next() (i uint64, ok bool)

	Mult() uint64
	Mod() uint64
	Start() uint64

	Min() uint64
	Max() uint64
}

// New creates a non-thread-safe iterator over the numeric range.
func New(mult, mod, min, max, start uint64) *iter {
	if start == 0 {
		start = rand.Uint64() % (max - min)
	}
	start = start % (max - min)
	if start == 0 {
		start = 1
	}

	return &iter{
		mult:  mult,
		mod:   mod,
		start: start,
		min:   min,
		max:   max,
	}
}

type iter struct {
	mult  uint64
	mod   uint64
	start uint64
	i     uint64
	max   uint64
	min   uint64
}

var _ Iter = (*iter)(nil)

func (i *iter) Mod() uint64   { return i.mod }
func (i *iter) Mult() uint64  { return i.mult }
func (i *iter) Start() uint64 { return i.start }
func (i *iter) Max() uint64   { return i.max }
func (i *iter) Min() uint64   { return i.min }

func (i *iter) Next() (v uint64, more bool) {
	v, more = i.next()
	for more && v > i.max {
		v, more = i.next()
	}
	return
}

func (i *iter) next() (v uint64, more bool) {
	if i.i == 0 {
		i.i = i.start
		more = true
	} else {
		i.i = (i.i * i.mult) % i.mod
		more = i.i != i.start
	}

	v = i.i + i.min - 1
	return
}
