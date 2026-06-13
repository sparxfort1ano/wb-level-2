// Package errors defines application-wide sentinel errors.
package errors

import "errors"

var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrBadArgument     = errors.New("bad argument")
)
