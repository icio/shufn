package shufn

// Calc determines the minimum mod and an appropriate mult for passing to New.
func Calc(rangeMin, rangeMax, start uint64) (mult, mod, min, max, start uint64) {
	if rangeMax < rangeMin {
		rangeMin, rangeMax = rangeMax, rangeMin
	}

	primes := primesPast(rangeMax - rangeMin + 1)
	if len(primes) == 0 {
		return 1, 2, rangeMin, rangeMax, start
	}

	mod = primes[len(primes)-1]
	roots := primePrimitiveRoots(mod, primes)
	mult = roots[len(roots)*2/3]

	return mult, mod, rangeMin, rangeMax, start
}

// primesPast returns all prime numbers up to the next prime above min.
func primesPast(min uint64) (primes []uint64) {
	if min < 2 {
		return nil
	}
Sieve:
	for p := uint64(2); ; p += 1 + (p & 1) {
		for _, prime := range primes {
			if p%prime == 0 {
				continue Sieve
			}
		}
		primes = append(primes, p)
		if p > min {
			return
		}
	}
}

// primePrimitiveRoots returns the primitive roots of the prime t.
// Based on https://math.stackexchange.com/a/133720.
func primePrimitiveRoots(t uint64, primes []uint64) []uint64 {
	// s is Euler's Totient phi(prime), where phi(t) := t - 1 if t is prime.
	s := t - 1
	factors := primeFactors(s, primes)

	// For any primitive root a, and any exponent p in the set of powers, a^p mod t != 1.
	powers := make([]uint64, len(factors))
	for i, factor := range factors {
		powers[i] = s / factor
	}

	// Find the lowest a >= 2 which satisfies a^p mod t != 1 for all p in powers.
	a, found := uint64(2), false
FirstRoot:
	for ; a < t; a++ {
		for _, power := range powers {
			if modExp(a, power, t) == 1 {
				continue FirstRoot
			}
		}

		found = true
		break
	}
	if !found {
		return nil
	}

	// Find the other primitive roots.
	roots := []uint64{a}
	for p := uint64(2); p < t; p++ {
		if gcd(s, p) == 1 {
			roots = append(roots, modExp(a, p, t))
		}
	}

	return roots
}

// primeFactors determines which primes are factors of n.
func primeFactors(n uint64, primes []uint64) []uint64 {
	factors := make([]uint64, 0)
	for _, f := range primes {
		if n == 1 {
			break
		}
		if n%f == 0 {
			factors = append(factors, f)
			n /= f
		}
		for n%f == 0 {
			n /= f
		}
	}
	return factors
}

// gcd calculates the greatest common divisor of a and b.
func gcd(a, b uint64) uint64 {
	if b > a {
		a, b = b, a
	}
	for {
		c := a % b
		if c == 0 {
			return b
		}
		a, b = b, c
	}
}

// modExp calculates modular exponentiation with exponentiation by squaring, O(log exponent).
// https://github.com/mgenware/go-modular-exponentiation/blob/73cf1e8d902ec78e27667be3e1b4d365de3c0aec/modexp.go#L31-L46
func modExp(base, exponent, modulus uint64) uint64 {
	if modulus == 1 {
		return 0
	}
	if exponent == 0 {
		return 1
	}

	result := modExp(base, exponent/2, modulus)
	result = (result * result) % modulus
	if exponent&1 != 0 {
		return ((base % modulus) * result) % modulus
	}
	return result % modulus
}
