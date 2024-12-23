package chain

import (
	"context"
	"net"

	"github.com/jxo-me/netx/core/chain"
	"github.com/jxo-me/netx/core/connector"
	"github.com/jxo-me/netx/core/dialer"
	net_dialer "github.com/jxo-me/netx/x/internal/net/dialer"
)

type Transport struct {
	dialer    dialer.IDialer
	connector connector.IConnector
	options   chain.TransportOptions
}

func NewTransport(d dialer.IDialer, c connector.IConnector, opts ...chain.TransportOption) *Transport {
	tr := &Transport{
		dialer:    d,
		connector: c,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(&tr.options)
		}
	}

	return tr
}

func (tr *Transport) Dial(ctx context.Context, addr string) (net.Conn, error) {
	netd := &net_dialer.Dialer{
		Interface: tr.options.IfceName,
		Netns:     tr.options.Netns,
	}
	if tr.options.SockOpts != nil {
		netd.Mark = tr.options.SockOpts.Mark
	}
	if tr.options.Route != nil && len(tr.options.Route.Nodes()) > 0 {
		netd.DialFunc = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return tr.options.Route.Dial(ctx, network, addr)
		}
	}
	opts := []dialer.DialOption{
		dialer.HostDialOption(tr.options.Addr),
		dialer.NetDialerDialOption(netd),
	}
	return tr.dialer.Dial(ctx, addr, opts...)
}

func (tr *Transport) Handshake(ctx context.Context, conn net.Conn) (net.Conn, error) {
	var err error
	if hs, ok := tr.dialer.(dialer.IHandshaker); ok {
		conn, err = hs.Handshake(ctx, conn,
			dialer.AddrHandshakeOption(tr.options.Addr))
		if err != nil {
			return nil, err
		}
	}
	if hs, ok := tr.connector.(connector.IHandshaker); ok {
		return hs.Handshake(ctx, conn)
	}
	return conn, nil
}

func (tr *Transport) Connect(ctx context.Context, conn net.Conn, network, address string) (net.Conn, error) {
	netd := &net_dialer.Dialer{
		Interface: tr.options.IfceName,
		Netns:     tr.options.Netns,
	}
	if tr.options.SockOpts != nil {
		netd.Mark = tr.options.SockOpts.Mark
	}
	return tr.connector.Connect(ctx, conn, network, address,
		connector.DialerConnectOption(netd),
	)
}

func (tr *Transport) Bind(ctx context.Context, conn net.Conn, network, address string, opts ...connector.BindOption) (net.Listener, error) {
	if binder, ok := tr.connector.(connector.IBinder); ok {
		return binder.Bind(ctx, conn, network, address, opts...)
	}
	return nil, connector.ErrBindUnsupported
}

func (tr *Transport) Multiplex() bool {
	if mux, ok := tr.dialer.(dialer.IMultiplexer); ok {
		return mux.Multiplex()
	}
	return false
}

func (tr *Transport) Options() *chain.TransportOptions {
	if tr != nil {
		return &tr.options
	}
	return nil
}

func (tr *Transport) Copy() chain.Transporter {
	tr2 := &Transport{}
	*tr2 = *tr
	return tr
}
