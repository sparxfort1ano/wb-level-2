// Package options is responsible for parsing the grep parameters.
package options

import (
	"errors"
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

// ArgsParsing parses grep parameters such as flags using `go-flags` functions
// and input data (file or os.Stdin, pattern).
// Returns error if flags parsing or file opening were failed or due to
// insufficient or excessive number of arguments.
func ArgsParsing() (*Options, error) {
	opts := &Options{}

	args, err := flags.Parse(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	if len(args) == 0 {
		return nil, fmt.Errorf("no search pattern")
	} else if len(args) == 1 {
		opts.Reader = os.Stdin
	} else if len(args) == 2 {
		file, err := os.Open(args[1])
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %w", args[1], err)
		}
		opts.Reader = file
	} else {
		return nil, fmt.Errorf("argument limit (2) exceeded")
	}

	opts.Pattern = args[0]

	return opts, nil
}

// ErrorHandling outputs to os.Stderr an error message and its code,
// depending on the type of error.
// It also ensures that the --help flag is handled correctly.
func ErrorHandling(err error) {
	var (
		flagsErr *flags.Error
	)

	switch {
	case errors.As(err, &flagsErr):
		if errors.Is(err, flags.ErrHelp) {
			os.Exit(0)
		}
	default:
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
