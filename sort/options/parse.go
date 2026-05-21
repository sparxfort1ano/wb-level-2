// Package options is responsible for parsing the sort parameters
// and containts sorting functions for them.
package options

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/jessevdk/go-flags"
)

// Close closes all the opened input files.
func (o *Options) Close() {
	for _, file := range o.Inputs {
		file.Close()
	}
}

// ArgsParsing parses sorting parameters such as flags using `go-flags` functions
// and input data (files or os.Stdin).
// Returns error if flags parsing or files opening were failed.
func ArgsParsing() (*Options, error) {
	opts := &Options{}

	files, err := flags.Parse(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	if len(files) == 0 {
		opts.Inputs = append(opts.Inputs, io.NopCloser(os.Stdin))
	} else {
		for _, path := range files {
			file, err := os.Open(path)
			if err != nil {
				opts.Close()
				return nil, fmt.Errorf("failed to open file %s: %w", path, err)
			}
			opts.Inputs = append(opts.Inputs, file)
		}
	}

	return opts, nil
}

// ErrorHandling outputs to os.Stderr an error message and its code,
// depending on the type of error.
// It also ensures that the --help flag is handled correctly.
func ErrorHandling(err error) {
	var (
		flagsErr *flags.Error
		pathErr  *os.PathError
	)

	switch {
	case errors.As(err, &flagsErr):
		if errors.Is(err, flags.ErrHelp) {
			os.Exit(0)
		}
	case errors.As(err, &pathErr):
		fmt.Fprintln(os.Stderr, err)

	default:
		fmt.Fprintf(os.Stderr, "Unknown error: %v\n", err)
	}
	os.Exit(1)
}
