package operator

import "strings"

// LogicParse splits a raw input line into a sequence of exectuable commands
// and logical operators.
func LogicParse(line string) []string {
	// The trick is used so as not to confuse `|` and `||`.
	line = strings.ReplaceAll(line, "&&", "\x00&&\x00")
	line = strings.ReplaceAll(line, "||", "\x00||\x00")

	rawTokens := strings.Split(line, "\x00")
	result := make([]string, 0, len(rawTokens))

	for _, token := range rawTokens {
		trimmed := strings.TrimSpace(token)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}
