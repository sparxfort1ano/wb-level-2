// Package cut contains all the working logic of the cut utility.
package cut

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/sparxfort1ano/wb-level-2/cut/options"
)

// RunCut serves as a core function for the cut utility.
// It receives an input stream and cut flags, handles it
// and prints the result to out.
func RunCut(out io.Writer, opts *options.Options) error {
	reader := bufio.NewReader(opts.Reader)

	for {
		line, err := reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return fmt.Errorf("failed to read string from an input file: %w", err)
		}

		line = strings.TrimRight(line, "\r\n")

		if line == "" && err != nil {
			break
		}

		var (
			countTokens    int
			token          strings.Builder
			firstTokenRead bool
		)

		isMatch := func() bool {
			return opts.SelectedFields[countTokens] || (opts.OpenEndedFrom != options.OpenEndedFromUninitialized && countTokens >= opts.OpenEndedFrom-1)
		}

		for _, symbol := range line {
			if string(symbol) == opts.Delimiter {
				if isMatch() {
					if firstTokenRead {
						fmt.Fprint(out, opts.Delimiter)
					}

					tokenStr := token.String()

					fmt.Fprint(out, tokenStr)

					token.Reset()
					firstTokenRead = true
				}
				countTokens++
				continue
			}

			if isMatch() {
				token.WriteRune(symbol)
			}
		}

		if countTokens == 0 {
			if opts.OnlySeparated {
				continue
			}

			fmt.Fprint(out, line)
		} else if isMatch() {
			if firstTokenRead {
				fmt.Fprint(out, opts.Delimiter)
			}

			tokenStr := token.String()

			fmt.Fprint(out, tokenStr)
		}

		fmt.Fprintln(out)
	}

	return nil
}
