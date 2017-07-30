package main

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/icio/shufn"
)

func main() {
	verbose := flag.Bool("v", false, "Print the iterator used to stderr.")
	sync := flag.Bool("t", false, "Use a thread-safe iterator (doesn't change logic).")
	quiet := flag.Bool("q", false, "Don't enumerate the iterator.")

	// Parse flags.
	flag.Usage = func() {
		usage(2, nil)
	}
	flag.Parse()

	// Determine the invocation style.
	args := flag.Args()
	minArg, startArg, maxArg, multArg, modArg := -1, -1, -1, -1, -1
	switch len(args) {
	case 1:
		maxArg = 0
	case 2:
		minArg, maxArg = 0, 1
	case 3:
		minArg, startArg, maxArg = 0, 1, 2
	case 5:
		minArg, startArg, maxArg, multArg, modArg = 0, 1, 2, 3, 4
	default:
		usage(1, nil)
	}

	// Read the arguments.
	min := parseArg(args, minArg, 1)
	start := parseArg(args, startArg, 0)
	max := parseArg(args, maxArg, 0)
	mult := parseArg(args, multArg, 0)
	mod := parseArg(args, modArg, 0)

	// Seed the random value which will be used for the starting position.
	rand.Seed(time.Now().UnixNano())

	// Construct the iterator.
	var it shufn.Iter
	if mult == 0 || mod == 0 {
		it = shufn.New(shufn.Calc(min, max, start))
	} else {
		it = shufn.New(mult, mod, min, max, start)
	}
	if *sync {
		it = shufn.Sync(it)
	}

	// Dump the iter config.
	if *verbose || *quiet {
		dump(os.Stderr, it, *sync)
	}

	// Enumerate the range.
	if !*quiet {
		for {
			i, ok := it.Next()
			if !ok {
				break
			}

			fmt.Println(i)
		}
	}
}

func parseArg(args []string, position int, def uint64) uint64 {
	if position == -1 {
		return def
	}
	arg, err := strconv.ParseUint(args[position], 10, 64)
	if err != nil {
		usage(1, err) // Exits.
	}
	return arg
}

func usage(exitCode int, err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
	fmt.Fprintln(os.Stderr, "Usage: shufn [-v|-q] [MIN [START]] MAX [MULT MOD]\n", os.Args[0])
	os.Exit(exitCode)
}

func dump(w io.Writer, it shufn.Iter, sync bool) {
	syncFix := func(s string) string { return s }
	if sync {
		syncFix = func(s string) string {
			return strings.Replace(strings.Replace(s, "\n", ")\n", 1), "= ", "= shufn.Sync(", 1)
		}
	}

	mult, mod, min, max, start := it.Mult(), it.Mod(), it.Min(), it.Max(), it.Start()
	fmt.Fprintf(w, "// Iterate from %d to %d, inclusive.\n", min, max)
	fmt.Fprintf(w, "// Repeat invocation: shufn %d %d %d\n", min, start, max)
	fmt.Fprintf(w, "// Instant invocation: shufn %d %d %d %d %d\n", min, start, max, mult, mod)
	fmt.Fprintf(w, "var iter        = %#v\n", it)
	fmt.Fprintf(w, syncFix("var constructor = shufn.New(%d, %d, %d, %d, %d)\n"), mult, mod, min, max, start)
	fmt.Fprintf(w, syncFix("var constrCalc  = shufn.New(shuf.Calc(%d, %d, %d))\n"), min, max, start)
	fmt.Fprintf(w, syncFix("var randRotated = shufn.New(%d, %d, %d, %d, 0)\n"), mult, mod, min, max)
	fmt.Fprintf(w, syncFix("var randRotCalc = shufn.New(shuf.Calc(%d, %d, 0))\n"), min, max)
}
