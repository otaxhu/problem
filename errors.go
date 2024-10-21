package problem

import "errors"

var (
	ErrInvalidContentType = errors.New("the Content-Type header is not one of 'application/problem+json' nor 'application/problem+xml'")
)
