package shufn_test

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/icio/shufn"
)

func ExampleSimple() {
	// it produces values in the range 0-100 inclusive. The start parameter is
	// passed through to New verbatim. When New receives start=0 it replaces it
	// with a random value.
	rand.Seed(0)
	it := shufn.New(shufn.Calc(1, 3, 0))
	for {
		i, more := it.Next()
		if !more {
			break
		}

		fmt.Println(i)
	}
	// Output: 1
	// 3
	// 2
}

func ExampleParallel() {
	work := make(chan int)
	done := make(chan int)

	// Distribute the random values via a channel.
	go func(it shufn.Iter) {
		for {
			i, more := it.Next()
			if !more {
				break
			}
			work <- int(i)
		}
		close(work)
	}(shufn.New(shufn.Calc(0, 10, 0)))

	// Create two workers to square each integer on the work queue and write to done.
	W, workers := 2, make(chan bool)
	for w := 0; w < W; w++ {
		go func(work chan int, done chan int) {
			for n := range work {
				done <- n * n
			}
			workers <- true
		}(work, done)
	}

	// Close the done channel after the workers have completed.
	go func() {
		for w := 0; w < W; w++ {
			<-workers
		}
		close(done)
	}()

	// Collect and sort the results.
	squares := make(chan []int)
	go func() {
		sqs := make([]int, 0, 10)
		for sq := range done {
			sqs = append(sqs, sq)
		}
		sort.Ints(sqs)
		squares <- sqs
	}()

	fmt.Println(<-squares)
	// Output: [0 1 4 9 16 25 36 49 64 81 100]
}
