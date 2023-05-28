package boot

import (
	"github.com/jxo-me/netx/core/connector"
	direct "github.com/jxo-me/netx/x/connector/direct"
	"github.com/jxo-me/netx/x/connector/forward"
	"github.com/jxo-me/netx/x/connector/http"
	"github.com/jxo-me/netx/x/connector/http2"
	"github.com/jxo-me/netx/x/connector/relay"
	"github.com/jxo-me/netx/x/connector/sni"
	v4 "github.com/jxo-me/netx/x/connector/socks/v4"
	v5 "github.com/jxo-me/netx/x/connector/socks/v5"
	"github.com/jxo-me/netx/x/connector/ss"
	ssu "github.com/jxo-me/netx/x/connector/ss/udp"
	"github.com/jxo-me/netx/x/connector/sshd"
	"github.com/jxo-me/netx/x/consts"
)

var Connectors = map[string]connector.NewConnector{
	consts.Direct:  direct.NewConnector,
	consts.Virtual: direct.NewConnector,
	consts.Forward: forward.NewConnector,
	consts.Http:    http.NewConnector,
	consts.Http2:   http2.NewConnector,
	consts.Relay:   relay.NewConnector,
	consts.Sni:     sni.NewConnector,
	consts.Socks4:  v4.NewConnector,
	consts.Socks4a: v4.NewConnector,
	consts.Socks5:  v5.NewConnector,
	consts.Socks:   v5.NewConnector,
	consts.Ss:      ss.NewConnector,
	consts.Ssu:     ssu.NewConnector,
	consts.Sshd:    sshd.NewConnector,
}
