package dialer

import (
	"context"
	"net"

	"github.com/jxo-me/netx/sdk/core/metadata"
)

// IDialer Transporter is responsible for dialing to the proxy server.
type IDialer interface {
	Init(metadata.IMetaData) error
	Dial(ctx context.Context, addr string, opts ...DialOption) (net.Conn, error)
}

type Handshaker interface {
	Handshake(ctx context.Context, conn net.Conn, opts ...HandshakeOption) (net.Conn, error)
}

type Multiplexer interface {
	Multiplex() bool
}
