package command

import (
	"fmt"
	"io"
	"strings"
)

// Echo prints the provided arguments to the specified output stream.
func Echo(outStream io.Writer, args []string) {
	str := strings.Join(args[1:], " ")

	if len(str) >= 2 {
		hasSingleQuotes := strings.HasPrefix(str, "'") && strings.HasSuffix(str, "'")
		hasDoubleQuotes := strings.HasPrefix(str, "\"") && strings.HasSuffix(str, "\"")
		if hasSingleQuotes || hasDoubleQuotes {
			str = str[1 : len(str)-1]
		}
	}

	fmt.Fprintln(outStream, str)
}
