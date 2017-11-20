package shufn

import (
	"errors"
	"math/rand"
)

var (
	ErrInvalidBounds = errors.New("shufn.Seq.Min,Max should define a positive range.")
	ErrInvalidMod    = errors.New("shufn.Seq.Mod should be greater than Max-Min, and prime.")
	ErrInvalidMult   = errors.New("shufn.Seq.Mult should be... something?")
	ErrInvalidLoop   = errors.New("Seq.Loop should be between Min and Max.")
)

type Seq struct {
	Min, Max  uint64
	Mult, Mod uint64
	Loop      uint64

	Primes       []uint64
	ModPrimRoots []uint64
}

type Iter struct {
	// Seq describes the sequence the Iter will follow.
	*Seq
	// I is the current value of the iterator.
	I uint64
	// isNext indicates whether the current I is actually the next one too.
	isNext bool
}

// MustIter returns an Iter for the given Seq, panicking should an error occur.
func MustIter(s *Seq) *Iter {
	iter, err := NewIter(s)
	if err != nil {
		panic(err)
	}
	return iter
}

// NewIter returns an Iter for the given Seq. If any of Mult, Mod, or Loop are
// not specified then these values are randomly chosen. If Mult or Mod need to
// be selected, Primes and ModPrimRoots on the Seq will be populated if empty.
// Use "math/rand".Seed to ensure a random sequence if desired.
func NewIter(s *Seq) (*Iter, error) {
	if err := validate(s); err != nil {
		return nil, err
	}

	if s.Mod == 0 || s.Mult == 0 {
		if s.Min == s.Max {
			// Special case where Min==Max.
			s.Mod = 2
			s.Mult = 1
		} else {
			// Determine possible values for Mod and Mult.
			if len(s.Primes) == 0 {
				s.Primes = primesPast(s.Max - s.Min + 1)
			}
			if s.Mod == 0 {
				if s.Max-s.Min == 0 {
					s.Mod = 1
				} else {
					s.Mod = s.Primes[len(s.Primes)-1]
				}
			}
			if len(s.ModPrimRoots) == 0 {
				s.ModPrimRoots = primePrimitiveRoots(s.Mod, s.Primes)
			}
		}
	}

	iter := &Iter{Seq: &(*s)}

	// Randomly choose which sequence to follow. There are only
	// len(s.ModPrimRoots) * (iter.Max - iter.Min) possible sequences.
	if iter.Mult == 0 {
		iter.Mult = s.ModPrimRoots[rand.Intn(len(s.ModPrimRoots))]
	}
	if iter.Loop == 0 {
		if iter.Min == iter.Max {
			iter.Loop = iter.Min
		} else {
			iter.Loop = rand.Uint64()%(iter.Max-iter.Min+1) + iter.Min
		}
		if iter.Loop == 0 {
			iter.Loop = 1
		}
	}

	// Kick-start the iterator.
	iter.I = iter.Loop
	iter.isNext = true

	return iter, nil
}

func validate(s *Seq) error {
	if s.Min > s.Max {
		return ErrInvalidBounds
	}
	if s.Mod != 0 && s.Mod <= s.Max-s.Min {
		return ErrInvalidMod
	}
	if s.Mult != 0 && s.Mult > s.Mod {
		return ErrInvalidMult
	}
	if s.Loop != 0 && (s.Loop < s.Min || s.Loop > s.Max) {
		return ErrInvalidLoop
	}
	return nil
}

// Next returns whether there are more numbers in the sequence; and indicates
// the next is available on i.I. (Not thread safe!)
func (i *Iter) Next() (more bool) {
	more = i.next()
	for more && i.I > i.Max {
		more = i.next()
	}
	return
}

func (i *Iter) NextI() (I uint64, more bool) {
	more, I = i.Next(), i.I
	return
}

func (i *Iter) next() bool {
	if i.isNext {
		// i.I already contains the next value.
		i.isNext = false
		return true
	}

	// i.I contains the previous value - lets find the next.
	i.I = ((i.I-i.Min+1)*i.Mult)%i.Mod + i.Min - 1
	return i.I != i.Loop
}

// Consume reads the remainder of iter into a []uint64.
func Consume(iter *Iter) []uint64 {
	s := make([]uint64, 0, iter.Max-iter.Min+1)
	for iter.Next() {
		s = append(s, iter.I)
	}
	return s
}
