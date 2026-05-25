package options

import "io"

// RawOptions represents raw cut options gotten from parsing flags using `go-flags` tags.
type RawOptions struct {
	FieldNumbers  string `short:"f" required:"true" description:"Specify the field (column) numbers to be printed"`
	Delimiter     string `short:"d" default:"\t" description:"Determinate a delimiter"`
	OnlySeparated bool   `short:"s" description:"Print only the lines, containing the delimiter"`
}

// Options represents full cut options converted from RawOptions and input data.
type Options struct {
	Reader         io.Reader
	Delimiter      string
	SelectedFields map[int]bool
	OpenEndedFrom  int
	OnlySeparated  bool
}

const (
	OpenEndedFromUninitialized = -1
)

// NewOptions creates a new instance of Options.
func NewOptions(
	delimiter string,
	onlySeparated bool,
	reader io.Reader,
) *Options {
	return &Options{
		SelectedFields: make(map[int]bool),
		Delimiter:      delimiter,
		OnlySeparated:  onlySeparated,
		Reader:         reader,
		OpenEndedFrom:  OpenEndedFromUninitialized,
	}
}
