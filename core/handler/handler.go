package handler

import (
	"context"
	"net"

	"github.com/jxo-me/netx/core/chain"
	"github.com/jxo-me/netx/core/metadata"
)

type IHandler interface {
	Init(metadata.IMetaData) error
	Handle(context.Context, net.Conn, ...HandleOption) error
}

type IForwarder interface {
	Forward(chain.IHop)
}
