package options

import "io"

// Options represents grep options for parsing: flags using `go-flags` tags
// and input data.
type Options struct {
	AfterContext    int  `short:"A" description:"Print N lines after the match"`
	BeforeContext   int  `short:"B" description:"Print N lines before the match"`
	AroundContext   int  `short:"C" description:"Print N lines after and before the match"`
	CountOnly       bool `short:"c" description:"Print number of matches"`
	IgnoreCase      bool `short:"i" description:"Case-insensitive search"`
	Reverse         bool `short:"v" description:"Print lines which dont match the pattern"`
	LiteralSearch   bool `short:"F" description:"Search an exact substring match"`
	ShowLineNumbers bool `short:"n" description:"Print lines which match the pattern showing their line number"`

	Pattern string
	Reader  io.Reader
}
