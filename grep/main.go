package main

import (
	"log"
	"os"

	"github.com/sparxfort1ano/wb-level-2/grep/grep"
	"github.com/sparxfort1ano/wb-level-2/grep/options"
)

func main() {
	opts, err := options.ArgsParsing()
	if err != nil {
		options.ErrorHandling(err)
	}

	out := os.Stdout
	if err := grep.RunGrep(out, opts); err != nil {
		log.Fatal(err)
	}
}
