// Package shufn implements a pseudo-random iterator over a numeric range with
// practically no memory overhead. shufn works on the principal that, given a
// prime number p, a primitive root r of p, then i' = (i * r) % p will loop
// through all integers [1, p) in a seemingly random order. This sequence is
// consistent upon invocations. The start parameter can be randomized to shift
// the sequence.
//
// shufn.Call provides a user-friendly interface for calculating an appropriate
// prime p and primitive root r for a given range [min, max]. shufn.New can
// then be used to instantiate the iterator immediately. Where the range is
// known at compilation time, the results of shufn.Call can be hard-coded.
//
// The shufn executable (go get github.com/icio/shufn/cmd/shufn) provides this
// precomputation on the command-line. For example: `shufn -q 1000000` will
// demonstrate multiple ways to use the shufn package.
//
// While it := shufn.New(...) is not thread-safe, shufn.Sync(it) is thread-safe
// and implements the same interface.
package shufn
