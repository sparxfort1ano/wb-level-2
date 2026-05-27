package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sparxfort1ano/wb-level-2/minishell/execute"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	streams := execute.NewStreams()

	for {
		fmt.Fprint(streams.OutputStream, "> ")

		if !scanner.Scan() {
			break
		}

		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		commands := strings.Split(line, "|")
		var err error
		switch len(commands) {
		case 1:
			err = streams.Execute(line)
		default:
			err = streams.ExecutePipes(commands)
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintln(streams.ErrorStream, err)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(streams.ErrorStream, err)
		os.Exit(1)
	}
}
