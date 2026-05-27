package operator

import "strings"

type RedirectOptions struct {
	CmdLine    string
	InputFile  string
	OutputFile string
}

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
