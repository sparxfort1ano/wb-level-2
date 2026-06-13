package request

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	errs "github.com/sparxfort1ano/wb-level-2/calendar/internal/core/errors"
)

// GetIntQueryParam extracts an integer query parameter
// from the HTTP request by its key.
// It returns nil if the parameter is missing
// or an error if the value is not a valid integer.
func GetIntQueryParam(r *http.Request, key string) (*int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	val, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key='%s' not a valid integer: %v: %w",
			param,
			key,
			err,
			errs.ErrInvalidArgument,
		)
	}

	return &val, nil
}

// GetDateQueryParam extracts a time.Time query parameter
// from the HTTP request by its key.
// It returns nil if the parameter is missing
// or an error if the value is not a valid time.Time variable.
func GetDateQueryParam(r *http.Request, key string) (*time.Time, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	layout := "2006-01-02"
	date, err := time.Parse(layout, param)
	if err != nil {
		return nil, fmt.Errorf(
			"param=`%s` by key=`%s` not a valid date: %v: %w",
			param,
			key,
			err,
			errs.ErrInvalidArgument,
		)
	}

	return &date, nil
}
