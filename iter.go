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

func (s *Seq) Validate() error {
	if s.Min >= s.Max {
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

type Iter struct {
	*Seq
	i, I uint64
}

func MustIter(s *Seq) *Iter {
	iter, err := NewIter(s)
	if err != nil {
		panic(err)
	}
	return iter
}

func NewIter(s *Seq) (*Iter, error) {
	// Prepare the sequence.
	if err := s.Validate(); err != nil {
		return nil, err
	}
	if s.Mod == 0 || s.Mult == 0 {
		// Calculate the primes and primitive roots required by our Iter.
		if len(s.Primes) == 0 {
			s.Primes = primesPast(s.Max - s.Min + 1)
		}
		if s.Mod == 0 {
			s.Mod = s.Primes[len(s.Primes)-1]
		}
		s.ModPrimRoots = primePrimitiveRoots(s.Mod, s.Primes)
	}

	iter := &Iter{Seq: &(*s)}

	// Choose the sequence if not all variables are given.
	if iter.Mult == 0 {
		iter.Mult = s.ModPrimRoots[rand.Intn(len(s.ModPrimRoots))]
	} else {
		// TODO: Validate Mult.
	}
	if iter.Loop == 0 {
		iter.Loop = rand.Uint64() % (iter.Max - iter.Min)
		if iter.Loop == 0 {
			iter.Loop = 1
		}
	} else {
		// TODO: Validate Loop.
	}
	iter.I = iter.Loop + iter.Min + 1

	return iter, nil
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
		// FIXME: This doesn't quite start correctly. Should be Loop?
		i.i = i.I - i.Min + 1
		more = true
	} else {
		i.i = (i.i * i.Mult) % i.Mod
		more = i.i != i.Loop
	}

	I = i.i + i.Min - 1
	return
}
