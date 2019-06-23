package main

import (
	"flag"
	"math/rand"

	"github.com/lovung/sortBigFile/lib/golib"
)

func main() {
	fo := flag.String("o", "resources/list.txt", "file path to write to")
	num := flag.Int("n", 1000000, "number of items")
	flag.Parse()
	var slice []int
	for i := 0; i < *num/100; i++ {
		slice = append(slice, rand.Perm(10000)...)
	}

	golib.WriteFile(*fo, slice, len(slice))
}
