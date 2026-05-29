package command

import (
	"fmt"
	"os"
)

var prevPWD string

// ChangeDirectory updates the current working directory of the shell process.
// It also supports `~`, `-` arguments.
func ChangeDirectory(args []string) error {
	var dir string
	if len(args) == 1 || (len(args) == 2 && args[1] == "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("$HOME not set: %w", err)
		}
		dir = homeDir
	} else if len(args) > 1 && args[1] == "-" {
		if prevPWD == "" {
			return fmt.Errorf("OLDPWD not set")
		}
		dir = prevPWD
	} else {
		dir = args[1]
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get path of working directory error: %w", err)
	}

	if err := os.Chdir(dir); err != nil {
		return fmt.Errorf("change directory error: %w", err)
	}

	prevPWD = currentDir

	return nil
}
