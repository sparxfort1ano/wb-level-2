package command

import (
	"fmt"
	"io"
	"os"
)

func PrintWorkingDirectory(outStream io.Writer) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get path of working directory error: %w", err)
	}

	fmt.Fprintln(outStream, currentDir)
	return nil
}
