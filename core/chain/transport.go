package chain

import (
	"context"
	"net"
	"time"

	net_dialer "github.com/jxo-me/netx/core/common/net/dialer"
	"github.com/jxo-me/netx/core/connector"
	"github.com/jxo-me/netx/core/dialer"
)

type TransportOptions struct {
	Addr     string
	IfceName string
	Netns    string
	SockOpts *SockOpts
	Route    IRoute
	Timeout  time.Duration
}

type TransportOption func(*TransportOptions)

func AddrTransportOption(addr string) TransportOption {
	return func(o *TransportOptions) {
		o.Addr = addr
	}
}

func InterfaceTransportOption(ifceName string) TransportOption {
	return func(o *TransportOptions) {
		o.IfceName = ifceName
	}
}

func NetnsTransportOption(netns string) TransportOption {
	return func(o *TransportOptions) {
		o.Netns = netns
	}
}

func SockOptsTransportOption(so *SockOpts) TransportOption {
	return func(o *TransportOptions) {
		o.SockOpts = so
	}
}

func RouteTransportOption(route IRoute) TransportOption {
	return func(o *TransportOptions) {
		o.Route = route
	}
}

func TimeoutTransportOption(timeout time.Duration) TransportOption {
	return func(o *TransportOptions) {
		o.Timeout = timeout
	}
}

type Transport struct {
	dialer    dialer.IDialer
	connector connector.IConnector
	options   TransportOptions
}

func NewTransport(d dialer.IDialer, c connector.IConnector, opts ...TransportOption) *Transport {
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
	netd := &net_dialer.NetDialer{
		Interface: tr.options.IfceName,
		Netns:     tr.options.Netns,
		Timeout:   tr.options.Timeout,
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
	netd := &net_dialer.NetDialer{
		Interface: tr.options.IfceName,
		Netns:     tr.options.Netns,
		Timeout:   tr.options.Timeout,
	}
	if tr.options.SockOpts != nil {
		netd.Mark = tr.options.SockOpts.Mark
	}
	return tr.connector.Connect(ctx, conn, network, address,
		connector.NetDialerConnectOption(netd),
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

func (tr *Transport) Options() *TransportOptions {
	if tr != nil {
		return &tr.options
	}
	return nil
}

func (tr *Transport) Copy() *Transport {
	tr2 := &Transport{}
	*tr2 = *tr
	return tr
}
