package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"

	"github.com/sparxfort1ano/wb-level-2/minishell/execute"
	"github.com/sparxfort1ano/wb-level-2/minishell/operator"
)

func main() {
	streams := execute.NewStreams(os.Stdin, os.Stdout, os.Stderr)
	scanner := bufio.NewScanner(streams.InputStream)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	go func() {
		for {
			<-sigChan
			fmt.Fprint(streams.OutputStream, "\n> ")
		}
	}()

	for {
		fmt.Fprint(streams.OutputStream, "> ")

		if !scanner.Scan() {
			fmt.Fprintln(streams.OutputStream)
			break
		}

		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		var (
			lastErr  error
			skipNext bool
		)

		tokens := operator.LogicParse(line)
		for _, token := range tokens {
			if token == "&&" {
				if lastErr != nil {
					skipNext = true
				}
				continue
			}

			if token == "||" {
				if lastErr == nil {
					skipNext = true
				}
				continue
			}

			if skipNext {
				skipNext = false
				continue
			}

			commands := strings.Split(token, "|")
			switch len(commands) {
			case 1:
				lastErr = streams.Execute(token)
			default:
				lastErr = streams.ExecutePipes(commands)
			}
		}

		if lastErr != nil {
			if lastErr == io.EOF {
				break
			}
			if len(tokens) == 1 {
				fmt.Fprintln(streams.ErrorStream, lastErr)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(streams.ErrorStream, err)
		os.Exit(1)
	}
}
