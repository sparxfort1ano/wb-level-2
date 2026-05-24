// Package grep contains all the working logic of the grep utility.
package grep

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/sparxfort1ano/wb-level-2/grep/options"
)

// RunGrep serves as a core function for the grep utility.
// It receives the pattern, input stream and grep flags, handles it
// and prints the result to out.
func RunGrep(out io.Writer, opts *options.Options) error {
	var isMatch func(line string) bool
	if opts.LiteralSearch {
		pattern := opts.Pattern

		if opts.IgnoreCase {
			pattern = strings.ToLower(pattern)
		}

		isMatch = func(line string) bool {
			if opts.IgnoreCase {
				line = strings.ToLower(line)
			}
			return strings.Contains(line, pattern)
		}
	} else {
		var ignoreCaseFlag string
		if opts.IgnoreCase {
			ignoreCaseFlag = "(?i)"
		}

		rePattern, err := regexp.Compile(ignoreCaseFlag + opts.Pattern)
		if err != nil {
			return fmt.Errorf("invalid regular expression syntax")
		}

		isMatch = func(line string) bool {
			return rePattern.MatchString(line)
		}
	}

	reader := bufio.NewReader(opts.Reader)
	var (
		countLines, countMatches int
		lastStrPrinted           int
	)
	linesAfterContext := max(opts.AfterContext, opts.AroundContext)
	linesBeforeContext := max(opts.BeforeContext, opts.AroundContext)
	afterContextIndex, beforeContextIndex := 0, -1
	beforeContextArr := make([]string, linesBeforeContext)

	for {
		countLines++
		line, err := reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return fmt.Errorf("failed to read string from an input file: %w", err)
		}

		line = strings.TrimRight(line, "\r\n")

		if line == "" && err != nil {
			break
		}

		matched := opts.Reverse != isMatch(line)
		if matched {
			if opts.CountOnly {
				countMatches++
			} else {
				if lastStrPrinted != 0 && (linesBeforeContext > 0 || linesAfterContext > 0) && (countLines-linesBeforeContext) > lastStrPrinted+1 {
					fmt.Fprintln(out, "--")
				}

				if linesBeforeContext > 0 {
					beforeContextIndex++
					if beforeContextIndex >= linesBeforeContext {
						for i := beforeContextIndex % linesBeforeContext; i < linesBeforeContext; i++ {
							fmt.Fprintln(out, beforeContextArr[i])
						}
					}
					for i := 0; i < beforeContextIndex%linesBeforeContext; i++ {
						fmt.Fprintln(out, beforeContextArr[i])
					}
				}

				if opts.ShowLineNumbers {
					line = fmt.Sprintf("%d:%s", countLines, line)
				}
				fmt.Fprintln(out, line)

				beforeContextIndex = -1
				afterContextIndex = linesAfterContext
				lastStrPrinted = countLines
			}
			continue
		}

		if afterContextIndex > 0 {
			if linesAfterContext > 0 {
				if opts.ShowLineNumbers {
					line = fmt.Sprintf("%d-%s", countLines, line)
				}
				fmt.Fprintln(out, line)

				afterContextIndex--
				lastStrPrinted = countLines
			}
		} else {
			if linesBeforeContext > 0 {
				beforeContextIndex++
				if opts.ShowLineNumbers {
					line = fmt.Sprintf("%d-%s", countLines, line)
				}
				beforeContextArr[beforeContextIndex%linesBeforeContext] = line
			}
		}

		if err != nil {
			break
		}
	}

	if opts.CountOnly {
		fmt.Fprintln(out, countMatches)
	}

	return nil
}
