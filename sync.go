package shufn

import (
	"sync"
)

// Sync returns a concurrency-safe wrapper around Iter.
func Sync(i Iter) *syncIter {
	return &syncIter{Iter: i}
}

type syncIter struct {
	Iter
	mu   sync.Mutex
	done bool
}

var _ Iter = (*syncIter)(nil)

func (s *syncIter) Next() (v uint64, more bool) {
	s.mu.Lock()
	if s.done {
		v, more = 0, false
	} else {
		v, more = s.Iter.Next()
	}
	s.mu.Unlock()
	return
}
