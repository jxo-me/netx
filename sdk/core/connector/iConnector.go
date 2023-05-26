package connector

import (
	"context"
	"net"

	"github.com/jxo-me/netx/sdk/core/metadata"
)

// IConnector is responsible for connecting to the destination address.
type IConnector interface {
	Init(metadata.IMetaData) error
	Connect(ctx context.Context, conn net.Conn, network, address string, opts ...ConnectOption) (net.Conn, error)
}

type Handshaker interface {
	Handshake(ctx context.Context, conn net.Conn) (net.Conn, error)
}
