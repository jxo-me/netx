package listener

import (
	"net"

	"github.com/jxo-me/netx/sdk/core/metadata"
)

// IListener is a server listener, just like a net.Listener.
type IListener interface {
	Init(metadata.IMetaData) error
	Accept() (net.Conn, error)
	Addr() net.Addr
	Close() error
}
