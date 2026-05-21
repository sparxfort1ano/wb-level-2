package main

import (
	"log"
	"os"

	"github.com/sparxfort1ano/wb-level-2/sort/options"
	"github.com/sparxfort1ano/wb-level-2/sort/sort"
)

func main() {
	opts, err := options.ArgsParsing()
	if err != nil {
		options.ErrorHandling(err)
	}

	out := os.Stdout
	if err := sort.RunSort(out, opts); err != nil {
		log.Fatal(err)
	}
}
