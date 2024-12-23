package boot

import (
	direct "github.com/jxo-me/netx/x/connector/direct"
	"github.com/jxo-me/netx/x/connector/forward"
	"github.com/jxo-me/netx/x/connector/http"
	"github.com/jxo-me/netx/x/connector/http2"
	"github.com/jxo-me/netx/x/connector/relay"
	"github.com/jxo-me/netx/x/connector/serial"
	"github.com/jxo-me/netx/x/connector/sni"
	v4 "github.com/jxo-me/netx/x/connector/socks/v4"
	v5 "github.com/jxo-me/netx/x/connector/socks/v5"
	"github.com/jxo-me/netx/x/connector/ss"
	ssu "github.com/jxo-me/netx/x/connector/ss/udp"
	"github.com/jxo-me/netx/x/connector/sshd"
	"github.com/jxo-me/netx/x/connector/tcp"
	"github.com/jxo-me/netx/x/connector/tunnel"
	"github.com/jxo-me/netx/x/connector/unix"
	"github.com/jxo-me/netx/x/consts"
	"github.com/jxo-me/netx/x/registry"
)

var Connectors = map[string]registry.NewConnector{
	consts.Direct:  direct.NewConnector,
	consts.Virtual: direct.NewConnector,
	consts.Forward: forward.NewConnector,
	consts.Http:    http.NewConnector,
	consts.Http2:   http2.NewConnector,
	consts.Relay:   relay.NewConnector,
	consts.Serial:  serial.NewConnector,
	consts.Sni:     sni.NewConnector,
	consts.Socks4:  v4.NewConnector,
	consts.Socks4a: v4.NewConnector,
	consts.Socks5:  v5.NewConnector,
	consts.Socks:   v5.NewConnector,
	consts.Ss:      ss.NewConnector,
	consts.Ssu:     ssu.NewConnector,
	consts.Sshd:    sshd.NewConnector,
	consts.Tcp:     tcp.NewConnector,
	consts.Tunnel:  tunnel.NewConnector,
	consts.Unix:    unix.NewConnector,
}
