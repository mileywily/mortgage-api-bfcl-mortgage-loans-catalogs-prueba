package http_fif

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

// DecodeBody decodes the JSON body of an HTTP response into the provided interface.
// It returns an error if the response or its body is nil, or if the decoding fails.
func DecodeBody(res *http.Response, v interface{}) error {
	if res == nil {
		return ErrResponseNil
	}
	if res.Body == nil {
		return ErrNilBody
	}
	return decodeReadCloser(res.Body, v)
}

// decodeReadCloser decodes the JSON body from an io.ReadCloser into the provided interface.
// It validates the decoded structure and returns an error if the decoding or validation fails.
func decodeReadCloser(rc io.ReadCloser, v interface{}) error {
	if err := json.NewDecoder(rc).Decode(v); err != nil {
		return fmt.Errorf("can't decode body into %T, %v: %w", v, err, ErrInvalidBody)
	}

	if err := validate.Struct(v); err != nil {
		var valid *validator.InvalidValidationError
		if !errors.As(err, &valid) {
			return fmt.Errorf("can't decode body %s: %w", err, ErrFailedValidation)
		}
	}
	return rc.Close()
}
