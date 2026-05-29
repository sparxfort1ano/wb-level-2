// Package command implements the execution logic for the minishell's
// built-in commands and provides a fallback mechanism for launching
// external system binaries.
package command

import (
	"fmt"
	"io"
	"os"
)

// PrintWorkingDirectory retrieves the absolute path of the current directory
// and writes it to the provided output stream.
func PrintWorkingDirectory(outStream io.Writer) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get path of working directory error: %w", err)
	}

	fmt.Fprintln(outStream, currentDir)
	return nil
}
