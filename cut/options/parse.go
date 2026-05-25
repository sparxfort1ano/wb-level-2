// Package options is responsible for parsing the cut parameters.
package options

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/jessevdk/go-flags"
)

// ArgsParsing parses cut parameters such as flags using `go-flags` functions
// and input data (file or os.Stdin).
// Returns error if flags parsing or file opening were failed or due to
// insufficient or excessive number of arguments.
func ArgsParsing() (*Options, error) {
	rawOpts := &RawOptions{}

	args, err := flags.Parse(rawOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	reader, err := getReader(args)
	if err != nil {
		return nil, err
	}

	options := NewOptions(rawOpts.Delimiter, rawOpts.OnlySeparated, reader)

	if err := options.selectedFields(rawOpts.FieldNumbers); err != nil {
		return nil, err
	}
	return options, nil
}

func getReader(args []string) (io.Reader, error) {
	if len(args) == 0 {
		return os.Stdin, nil
	} else if len(args) == 1 {
		file, err := os.Open(args[0])
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %w", args[0], err)
		}
		return file, nil
	} else {
		return nil, fmt.Errorf("error: argument limit (1) exceeded")
	}
}

func (o *Options) selectedFields(fieldNumbers string) error {
	fieldsParsed := strings.Split(fieldNumbers, ",")
	for _, field := range fieldsParsed {
		splitByDash := strings.Split(field, "-")

		if len(splitByDash) == 1 {
			col, err := getPositiveIntCol(splitByDash[0])
			if err != nil {
				return err
			}
			o.SelectedFields[col-1] = true

		} else if len(splitByDash) == 2 {
			if splitByDash[0] == "" && splitByDash[1] == "" {
				return fmt.Errorf("error: invalid range with no endpoint: -")
			}

			var col1 int
			if splitByDash[0] == "" {
				col1 = 1
			} else {
				var err error
				col1, err = getPositiveIntCol(splitByDash[0])
				if err != nil {
					return err
				}
			}

			var col2 int
			if splitByDash[1] == "" {
				if o.OpenEndedFrom == OpenEndedFromUninitialized || col1 < o.OpenEndedFrom {
					o.OpenEndedFrom = col1
				}
				continue
			} else {
				var err error
				col2, err = getPositiveIntCol(splitByDash[1])
				if err != nil {
					return err
				}
			}

			if col1 > col2 {
				return fmt.Errorf("error: invalid decreasing range")
			}

			for col := col1; col <= col2; col++ {
				o.SelectedFields[col-1] = true
			}

		} else {
			return fmt.Errorf("error: invalid range specification")
		}
	}

	return nil
}

func getPositiveIntCol(splitByDashElem string) (int, error) {
	col, err := strconv.Atoi(splitByDashElem)
	if err != nil {
		return 0, fmt.Errorf("failed to convert string to integer: %w", err)
	}
	if col < 1 {
		return 0, fmt.Errorf("error: fields are numbered from 1")
	}

	return col, nil
}

// ErrorHandling outputs to os.Stderr an error message and its code,
// depending on the type of error.
// It also ensures that the --help flag is handled correctly.
func ErrorHandling(err error) {
	var flagsErr *flags.Error

	switch {
	case errors.As(err, &flagsErr):
		if flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}

		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	default:
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
