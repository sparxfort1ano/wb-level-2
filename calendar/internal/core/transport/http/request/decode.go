// Package request provides utilities for parsing incoming HTTP requests.
// It handles JSON decoding and structural or functional validation of incoming payloads.
package request

import (
	"fmt"
	"net/http"

	"github.com/go-playground/form"
	"github.com/go-playground/validator/v10"
	errs "github.com/sparxfort1ano/wb-level-2/calendar/internal/core/errors"
)

var requestValidator = validator.New()

// DecodeAndValidateRequest decodes the URL body of an HTTP request into the provided
// destination struct and validates its fields bases on struct tags.
// It returns an ErrInvalidArgument if decoding or validation fails.
func DecodeAndValidateRequest(r *http.Request, dest any) error {
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf(
			"parse url: %v: %w",
			err,
			errs.ErrInvalidArgument,
		)
	}

	decoder := form.NewDecoder()
	if err := decoder.Decode(&dest, r.PostForm); err != nil {
		return fmt.Errorf(
			"decode url: %v: %w",
			err,
			errs.ErrInvalidArgument,
		)
	}

	if err := requestValidator.Struct(dest); err != nil {
		return fmt.Errorf(
			"request validation: %v: %w",
			err,
			errs.ErrInvalidArgument,
		)
	}

	return nil
}
