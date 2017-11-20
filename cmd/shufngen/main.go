package main

import (
	. "github.com/icio/shufn"

	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())

	iter := MustIter(NewRange(5, 100000))
	for iter.Next() {
		fmt.Println(iter.I)
	}
}
