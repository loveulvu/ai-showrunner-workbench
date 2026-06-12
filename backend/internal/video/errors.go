package video

import "errors"

type ErrorKind string

const (
	ErrorKindRequest  ErrorKind = "request"
	ErrorKindUpstream ErrorKind = "upstream"
)

type Error struct {
	Kind      ErrorKind
	Message   string
	Err       error
	Retryable bool
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

func IsRequestError(err error) bool {
	var videoErr *Error
	return errors.As(err, &videoErr) && videoErr.Kind == ErrorKindRequest
}

func IsUpstreamError(err error) bool {
	var videoErr *Error
	return errors.As(err, &videoErr) && videoErr.Kind == ErrorKindUpstream
}
