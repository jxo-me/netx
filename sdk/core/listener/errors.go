package listener

import "errors"

var (
	ErrClosed = errors.New("accpet on closed listener")
)

type AcceptError struct {
	err error
}

func NewAcceptError(err error) error {
	return &AcceptError{err: err}
}

func (e *AcceptError) Error() string {
	return e.err.Error()
}

func (e *AcceptError) Timeout() bool {
	return false
}

func (e *AcceptError) Temporary() bool {
	return true
}

func (e *AcceptError) Unwrap() error {
	return e.err
}
