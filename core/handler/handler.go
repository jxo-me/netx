package handler

import (
	"context"
	"net"

	"github.com/jxo-me/netx/core/chain"
	"github.com/jxo-me/netx/core/metadata"
)

type Handler interface {
	Init(metadata.Metadata) error
	Handle(context.Context, net.Conn, ...HandleOption) error
}

type Forwarder interface {
	Forward(chain.Hop)
}
