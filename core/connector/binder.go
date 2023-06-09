package connector

import (
	"context"
	"errors"
	"net"
)

var (
	ErrBindUnsupported = errors.New("bind unsupported")
)

type IBinder interface {
	Bind(ctx context.Context, conn net.Conn, network, address string, opts ...BindOption) (net.Listener, error)
}
