// Package operator provides parsing utilities for logical operators
// and i/o redirections.
package operator

import "strings"

// RedirectOptions holds the parsed file paths for input and output
// redirections along with the trimmed command string ready for execution.
type RedirectOptions struct {
	CmdLine    string
	InputFile  string
	OutputFile string
}

// NewRedirectOptions creates a new instance of RedirectOptions.
func NewRedirectOptions(
	cmdLine string,
	inputFile string,
	outputFile string,
) *RedirectOptions {
	return &RedirectOptions{
		CmdLine:    cmdLine,
		InputFile:  inputFile,
		OutputFile: outputFile,
	}
}

// ParseRedirect extracts input and output file redirections from
// the given command line.
func ParseRedirect(line string) *RedirectOptions {
	var (
		cmdLine, outputFile, inputFile string
	)

	if len(line) == 0 {
		return NewRedirectOptions(cmdLine, inputFile, outputFile)
	}

	if idx := strings.Index(line, ">"); idx != -1 {
		tmp := strings.Split(line, ">")
		cmdLine = strings.TrimSpace(tmp[0])
		outputFile = strings.TrimSpace(tmp[1])
	}

	if idx := strings.Index(line, "<"); idx != -1 {
		tmp := strings.Split(line, "<")
		cmdLine = strings.TrimSpace(tmp[0])
		inputFile = strings.TrimSpace(tmp[1])
	}

	return NewRedirectOptions(cmdLine, inputFile, outputFile)
}
