package main

import (
	"log"
	"os"

	"github.com/sparxfort1ano/wb-level-2/cut/cut"
	"github.com/sparxfort1ano/wb-level-2/cut/options"
)

func main() {
	opts, err := options.ArgsParsing()
	if err != nil {
		options.ErrorHandling(err)
	}

	out := os.Stdout
	if err := cut.RunCut(out, opts); err != nil {
		log.Fatal(err)
	}
}
