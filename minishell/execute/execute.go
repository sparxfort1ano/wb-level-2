// Package execute provides the core execution engine for minishell.
// It handles i/o stream management, environment variable expansion,
// file redirections and commands routing (external, internal).
package execute

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/sparxfort1ano/wb-level-2/minishell/command"
	"github.com/sparxfort1ano/wb-level-2/minishell/operator"
)

// Streams encapsulates the default shell input, output and error streams.
type Streams struct {
	InputStream  io.Reader
	OutputStream io.Writer
	ErrorStream  io.Writer
}

// NewStreams creates a new instance of Streams.
func NewStreams(
	inputStream io.Reader,
	outputStream io.Writer,
	errorStream io.Writer,
) *Streams {
	return &Streams{
		InputStream:  inputStream,
		OutputStream: outputStream,
		ErrorStream:  errorStream,
	}
}

// executeNode processes and runs a single command within the shell.
// It expands environment variables actions, takes i/o redirection parsing,
// routes the command to appropriate executor (internal, external).
func (s *Streams) executeNode(line string, defaultOut io.Writer, defaultIn io.Reader) error {
	line = os.ExpandEnv(line)

	redirOpts := operator.ParseRedirect(line)

	inStream := defaultIn
	outStream := defaultOut
	errStream := s.ErrorStream

	if redirOpts.InputFile != "" {
		file, err := os.Open(redirOpts.InputFile)
		if err != nil {
			return fmt.Errorf("open file error")
		}
		defer file.Close()
		inStream = file
	}

	if redirOpts.OutputFile != "" {
		file, err := os.Create(redirOpts.OutputFile)
		if err != nil {
			return fmt.Errorf("create file error")
		}
		defer file.Close()
		outStream = file
	}

	if redirOpts.CmdLine == "" {
		redirOpts.CmdLine = line
	}

	args := strings.Fields(redirOpts.CmdLine)
	if len(args) == 0 {
		return nil
	}

	if strings.Contains(args[0], "=") {
		environ := strings.SplitN(args[0], "=", 2)
		os.Setenv(environ[0], environ[1])
		return nil
	}

	var err error
	switch args[0] {
	case "exit":
		err = io.EOF
	case "cd":
		err = command.ChangeDirectory(args)
	case "pwd":
		err = command.PrintWorkingDirectory(outStream)
	case "echo":
		command.Echo(outStream, args)
	case "ps":
		err = command.ProcessStatus(outStream)
	case "kill":
		err = command.Kill(args)
	default:
		err = command.ExecuteExternalCommand(outStream, inStream, errStream, args)
	}

	return err
}

// Execute runs a single, standalone command string using the default shell streams.
func (s *Streams) Execute(line string) error {
	return s.executeNode(line, s.OutputStream, s.InputStream)
}

// ExecutePipes orchestrates the concurrent execution of multiple commands
// connected by pipelines.
func (s *Streams) ExecutePipes(commands []string) error {
	n := len(commands)
	ins := make([]io.Reader, n)
	outs := make([]io.Writer, n)

	ins[0] = s.InputStream
	outs[n-1] = s.OutputStream

	for i := range n - 1 {
		r, w, err := os.Pipe()
		if err != nil {
			return fmt.Errorf("pipe error: %w", err)
		}

		ins[i+1] = r
		outs[i] = w
	}

	var (
		errs error
		wg   sync.WaitGroup
		mu   sync.Mutex
	)

	wg.Add(n)
	for i, command := range commands {
		go func(cmd string, w io.Writer, r io.Reader) {
			defer wg.Done()

			if err := s.executeNode(cmd, w, r); err != nil {
				mu.Lock()
				errs = errors.Join(errs, err)
				mu.Unlock()
			}

			if closer, ok := w.(io.WriteCloser); ok && w != s.OutputStream {
				closer.Close()
			}

			if closer, ok := r.(io.ReadCloser); ok && r != s.InputStream {
				closer.Close()
			}
		}(command, outs[i], ins[i])
	}

	wg.Wait()

	return errs
}
