package listener

import (
	"errors"
	"net"

	"github.com/jxo-me/netx/core/metadata"
)

var (
	ErrClosed = errors.New("accept on closed listener")
)

// IListener is a server listener, just like a net.Listener.
type IListener interface {
	Init(metadata.IMetaData) error
	Accept() (net.Conn, error)
	Addr() net.Addr
	Close() error
}

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

type NewListener func(opts ...Option) IListener
