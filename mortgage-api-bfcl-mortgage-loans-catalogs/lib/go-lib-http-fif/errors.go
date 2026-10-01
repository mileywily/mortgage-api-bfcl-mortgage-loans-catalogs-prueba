package http_fif

import "errors"

var (
	// ErrInvalidURLAddress is returned when the provided URL address is invalid.
	ErrInvalidURLAddress = errors.New("the provided URL address is invalid")

	// ErrRequestNil is returned when the request object is nil.
	ErrRequestNil = errors.New("the request object is nil")

	// ErrResponseNil is returned when the response object is nil.
	ErrResponseNil = errors.New("the response object is nil")

	// ErrFailedValidation is returned when the validation process fails.
	ErrFailedValidation = errors.New("the validation process failed")

	// ErrNilBody is returned when the body of the request or response is nil.
	ErrNilBody = errors.New("the body of the request or response is nil")

	// ErrInvalidBody is returned when the body of the request or response is invalid.
	ErrInvalidBody = errors.New("the body of the request or response is invalid")

	// ErrInvalidParametersCount is returned when the number of parameters provided is incorrect.
	ErrInvalidParametersCount = errors.New("the number of parameters provided is incorrect")
)
