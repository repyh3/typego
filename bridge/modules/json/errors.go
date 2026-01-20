package json

import "errors"

var (
	// ErrMissingArgument is returned when a required argument is missing
	ErrMissingArgument = errors.New("missing required argument")
)
