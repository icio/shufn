// Package shufn implements a pseudo-random iterator over a numeric range with
// practically no memory overhead.
//
// shufn works on the principal that, given a prime number Mod, a primitive
// root Mult of Mod, then i' = (i * Mult) % Mod will loop through all integers
// [1, Mod) in a seemingly random order. The sequence of random numbers is
// consistent for any (Mult, Mod) pair, but can be rotated according to Start.
//
// To calculate Mod and Mult for your desired range, see shufn.Calc.
//
// The shufn executable (go get github.com/icio/shufn/cmd/shufn) provides this
// precomputation on the command-line. For example: `shufn -q 1000000` will
// demonstrate multiple ways to use the shufn package.
//
// Note that the iterator is not thread-safe.
package shufn
