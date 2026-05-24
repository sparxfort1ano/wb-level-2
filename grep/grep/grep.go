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
	var pattern any
	if opts.LiteralSearch {
		pattern = opts.Pattern
		if opts.IgnoreCase {
			pattern = strings.ToLower(pattern.(string))
		}
	} else {
		var err error
		var ignoreCaseFlag string
		if opts.IgnoreCase {
			ignoreCaseFlag = "(?i)"
		}
		pattern, err = regexp.Compile(ignoreCaseFlag + opts.Pattern)
		if err != nil {
			return fmt.Errorf("invalid regular expression syntax")
		}
	}

	reader := bufio.NewReader(opts.Reader)
	var countLines, countMatches int
	var lastStrPrinted int
	linesAfterContext := max(opts.AfterContext, opts.AroundContext)
	linesBeforeContext := max(opts.BeforeContext, opts.AroundContext)
	afterContextIndex, beforeContextIndex := linesAfterContext, -1
	beforeContextArr := make([]string, linesBeforeContext)
	for {
		countLines++
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				if line != "" {
					if !strings.HasSuffix(line, "\n") {
						line += "\n"
					}
					var matched bool
					if opts.Reverse {
						switch pattern := pattern.(type) {
						case string:
							if opts.IgnoreCase {
								pattern = strings.ToLower(pattern)
							}
							if !strings.Contains(line, pattern) {
								matched = true
							}
						case *regexp.Regexp:
							if !pattern.MatchString(line) {
								matched = true
							}
						default:
							return fmt.Errorf("invalid pattern type")
						}
					} else {
						switch pattern := pattern.(type) {
						case string:
							if opts.IgnoreCase {
								pattern = strings.ToLower(pattern)
							}
							if strings.Contains(line, pattern) {
								matched = true
							}
						case *regexp.Regexp:
							if pattern.MatchString(line) {
								matched = true
							}
						default:
							return fmt.Errorf("invalid pattern type")
						}
					}

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
										fmt.Fprint(out, beforeContextArr[i])
									}
								}
								for i := 0; i < beforeContextIndex%linesBeforeContext; i++ {
									fmt.Fprint(out, beforeContextArr[i])
								}
							}

							if opts.ShowLineNumbers {
								line = fmt.Sprintf("%d:%s", countLines, line)
							}
							fmt.Fprint(out, line)
						}
					}

					if opts.CountOnly {
						fmt.Fprintln(out, countMatches)
					}
				}

				return nil
			}
			return fmt.Errorf("failed to read string from an input file: %w", err)
		}

		var matched bool
		if opts.Reverse {
			switch pattern := pattern.(type) {
			case string:
				if opts.IgnoreCase {
					pattern = strings.ToLower(pattern)
				}
				if !strings.Contains(line, pattern) {
					matched = true
				}
			case *regexp.Regexp:
				if !pattern.MatchString(line) {
					matched = true
				}
			default:
				return fmt.Errorf("invalid pattern type")
			}
		} else {
			switch pattern := pattern.(type) {
			case string:
				if opts.IgnoreCase {
					pattern = strings.ToLower(pattern)
				}
				if strings.Contains(line, pattern) {
					matched = true
				}
			case *regexp.Regexp:
				if pattern.MatchString(line) {
					matched = true
				}
			default:
				return fmt.Errorf("invalid pattern type")
			}
		}

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
							fmt.Fprint(out, beforeContextArr[i])
						}
					}
					for i := 0; i < beforeContextIndex%linesBeforeContext; i++ {
						fmt.Fprint(out, beforeContextArr[i])
					}
				}

				if opts.ShowLineNumbers {
					line = fmt.Sprintf("%d:%s", countLines, line)
				}
				fmt.Fprint(out, line)

				beforeContextIndex = -1 - linesAfterContext
				afterContextIndex = 0
				lastStrPrinted = countLines
			}
			continue
		}

		beforeContextIndex++
		if linesBeforeContext > 0 && beforeContextIndex >= 0 {
			if opts.ShowLineNumbers {
				line = fmt.Sprintf("%d-%s", countLines, line)
			}
			beforeContextArr[beforeContextIndex%linesBeforeContext] = line
		}

		if linesAfterContext > 0 && afterContextIndex < linesAfterContext {
			if opts.ShowLineNumbers {
				line = fmt.Sprintf("%d-%s", countLines, line)
			}
			fmt.Fprint(out, line)

			afterContextIndex++
			lastStrPrinted = countLines
		}
	}
}

// "error", -B 2 -A 2
// ijwiejw4io3n4oi3nw -> 0 0
// 4j23poek43pow,mpl32,e  -> 1 1
// ij34iiji43j34ij34o -> 0 2
// errorrrr при 0 2 if beforeContextIndex >= lines for (beforeContextIndex % lines) to n, 0 to (beforeContextIndex % lines)
// beforeContextIndex = -1 - linesAfterContext
// afterContextIndex = 0

// switch from after to before
// before after
// -2 1 print after
// -1 2 print after
// 0 2 put before into array

// found match while searching things after
// before after
// -2 1 print after
// BOOM match
// -2 1 print after

// -- printing
// before after
// -2 1 print after  line number x
// -2 2 print after  line number x+1 lastStrPrinted = x+1
// 0 2 to array (needa put here --) line number x+2
// 1 2 to array line number x+3
// 2 2 to array (will be printed) line number x+4
// 3 2 (will be printed) line number x+5
// BOOM match line number x+6

// criteria: x+6 - (linesBeforeContext=2) = x+4 > lastStrPrinted (x+1) + 1
// if true => --
