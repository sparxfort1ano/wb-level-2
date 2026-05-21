package options

import (
	"io"
)

// Options represents sorting options for parsing: flags using `go-flags` tags
// and input data.
type Options struct {
	ByColumn             int  `short:"k" description:"Sort by column N"`
	ByValue              bool `short:"n" description:"Sort by numerical value"`
	ByMonth              bool `short:"M" description:"Sort by months (Jan, Feb, etc)"`
	BySize               bool `short:"H" description:"Sort by size in bytes"`
	Reverse              bool `short:"r" description:"Sort by reverse"`
	Unique               bool `short:"u" description:"Sort unique strings"`
	IgnoreTrailingBlanks bool `short:"b" description:"Sort using ignoring trailing blanks"`
	IsSorted             bool `short:"c" description:"Check if the data is sorted"`

	Inputs []io.ReadCloser
}
