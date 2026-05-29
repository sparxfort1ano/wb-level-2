package command

import (
	"fmt"
	"os"
	"strconv"
)

// Kill sends a termination signal to a system process identified by its PID.
// Returns error if the PID is invalid or missing.
func Kill(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("only PID must be specified")
	}

	pid, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("pid must be integer: %w", err)
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("pid not found: %w", err)
	}

	if err := proc.Kill(); err != nil {
		return fmt.Errorf("failed to kill process with pid %d: %w", pid, err)
	}

	return nil
}
