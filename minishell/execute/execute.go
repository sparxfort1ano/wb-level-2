package execute

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/sparxfort1ano/wb-level-2/minishell/command"
	"github.com/sparxfort1ano/wb-level-2/minishell/operator"
)

type Streams struct {
	InputStream  io.Reader
	OutputStream io.Writer
	ErrorStream  io.Writer
}

func NewStreams() *Streams {
	return &Streams{
		InputStream:  os.Stdin,
		OutputStream: os.Stdout,
		ErrorStream:  os.Stderr,
	}
}

func (s *Streams) executeExternalCommand(args []string, redirOpts *operator.RedirectOptions) error {
	path, err := exec.LookPath(args[0])
	if err != nil {
		return fmt.Errorf("command not found: %w", err)
	}

	cmd := exec.Command(path, args[1:]...)

	if redirOpts.InputFile != "" {
		file, err := os.Open(redirOpts.InputFile)
		if err != nil {
			return fmt.Errorf("open file error")
		}
		defer file.Close()
		cmd.Stdin = file
	} else {
		cmd.Stdin = s.InputStream
	}

	if redirOpts.OutputFile != "" {
		file, err := os.Create(redirOpts.OutputFile)
		if err != nil {
			return fmt.Errorf("create file error")
		}
		defer file.Close()
		cmd.Stdout = file
	} else {
		cmd.Stdout = s.OutputStream
	}
	cmd.Stderr = s.ErrorStream

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command run error: %w", err)
	}

	return nil
}

func (s *Streams) Execute(line string) error {
	redirOpts := operator.ParseRedirect(line)
	if redirOpts.CmdLine == "" {
		redirOpts.CmdLine = line
	}
	args := strings.Fields(redirOpts.CmdLine)

	var err error
	switch args[0] {
	case "exit":
		err = io.EOF
	case "cd":
		err = command.ChangeDirectory(args)
	case "pwd":
		err = command.PrintWorkingDirectory(s.OutputStream)
	// TODO echo, kill, ps
	default:
		err = s.executeExternalCommand(args, redirOpts)
	}

	return err
}

func (s *Streams) ExecutePipes(commands []string) error {
	var cmds []*exec.Cmd
	for _, cmdStr := range commands {
		args := strings.Fields(cmdStr)

		path, err := exec.LookPath(args[0])
		if err != nil {
			return fmt.Errorf("command not found: %w", err)
		}
		cmd := exec.Command(path, args[1:]...)
		cmds = append(cmds, cmd)
	}

	var pipes []*os.File
	for i := range len(cmds) - 1 {
		r, w, err := os.Pipe()
		if err != nil {
			return fmt.Errorf("pipe error: %w", err)
		}
		cmds[i+1].Stdin = r
		cmds[i].Stdout = w
		pipes = append(pipes, r, w)
	}

	cmds[0].Stdin = s.InputStream
	cmds[len(cmds)-1].Stdout = s.OutputStream
	cmds[len(cmds)-1].Stderr = s.ErrorStream

	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("pipe command execution error: %w", err)
		}
	}

	for i := range len(cmds) - 1 {
		pipes[i*2+1].Close()
	}

	for i := len(cmds) - 1; i >= 0; i-- {
		if err := cmds[i].Wait(); err != nil {
			return fmt.Errorf("command wait error: %w", err)
		}
	}

	return nil
}
