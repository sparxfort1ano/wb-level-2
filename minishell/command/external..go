package command

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
)

// ExecuteExternalCommand locates and executes a system binary based on 
// the provided arguments.
// It listens for the SIGINT to gracefully terminate the execution without
// stopping the main shell instance.
func ExecuteExternalCommand(outStream io.Writer, inStream io.Reader, errStream io.Writer, args []string) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	path, err := exec.LookPath(args[0])
	if err != nil {
		return fmt.Errorf("command not found: %w", err)
	}

	cmd := exec.CommandContext(ctx, path, args[1:]...)

	cmd.Stdin = inStream
	cmd.Stdout = outStream
	cmd.Stderr = errStream

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command run error: %w", err)
	}

	return nil
}
