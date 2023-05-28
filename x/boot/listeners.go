package boot

import (
	"github.com/jxo-me/netx/core/listener"
	"github.com/jxo-me/netx/x/consts"
	listenerDns "github.com/jxo-me/netx/x/listener/dns"
	listenerDtls "github.com/jxo-me/netx/x/listener/dtls"
	listenerFtcp "github.com/jxo-me/netx/x/listener/ftcp"
	listenerGrpc "github.com/jxo-me/netx/x/listener/grpc"
	listenerHttp2 "github.com/jxo-me/netx/x/listener/http2"
	listenerHttpH2 "github.com/jxo-me/netx/x/listener/http2/h2"
	listenerHttp3 "github.com/jxo-me/netx/x/listener/http3"
	listenerHttpH3 "github.com/jxo-me/netx/x/listener/http3/h3"
	listenerIcmp "github.com/jxo-me/netx/x/listener/icmp"
	listenerKcp "github.com/jxo-me/netx/x/listener/kcp"
	listenerMtls "github.com/jxo-me/netx/x/listener/mtls"
	listenerMws "github.com/jxo-me/netx/x/listener/mws"
	listenerObfsHttp "github.com/jxo-me/netx/x/listener/obfs/http"
	listenerObfsTls "github.com/jxo-me/netx/x/listener/obfs/tls"
	listenerPht "github.com/jxo-me/netx/x/listener/pht"
	listenerQuic "github.com/jxo-me/netx/x/listener/quic"
	listenerRedirectTcp "github.com/jxo-me/netx/x/listener/redirect/tcp"
	listenerRedirectUdp "github.com/jxo-me/netx/x/listener/redirect/udp"
	listenerRtcp "github.com/jxo-me/netx/x/listener/rtcp"
	listenerRudp "github.com/jxo-me/netx/x/listener/rudp"
	listenerSsh "github.com/jxo-me/netx/x/listener/ssh"
	listenerSshd "github.com/jxo-me/netx/x/listener/sshd"
	listenerTap "github.com/jxo-me/netx/x/listener/tap"
	listenerTcp "github.com/jxo-me/netx/x/listener/tcp"
	listenerTls "github.com/jxo-me/netx/x/listener/tls"
	listenerTun "github.com/jxo-me/netx/x/listener/tun"
	listenerUdp "github.com/jxo-me/netx/x/listener/udp"
	listenerWs "github.com/jxo-me/netx/x/listener/ws"
)

var Listeners = map[string]listener.NewListener{
	consts.Dns:      listenerDns.NewListener,
	consts.Dtls:     listenerDtls.NewListener,
	consts.Ftcp:     listenerFtcp.NewListener,
	consts.Grpc:     listenerGrpc.NewListener,
	consts.Http2:    listenerHttp2.NewListener,
	consts.H2c:      listenerHttpH2.NewListener,
	consts.H2:       listenerHttpH2.NewTLSListener,
	consts.Http3:    listenerHttp3.NewListener,
	consts.H3:       listenerHttpH3.NewListener,
	consts.Icmp:     listenerIcmp.NewListener,
	consts.Kcp:      listenerKcp.NewListener,
	consts.Mtls:     listenerMtls.NewListener,
	consts.Mws:      listenerMws.NewListener,
	consts.Mwss:     listenerMws.NewTLSListener,
	consts.Ohttp:    listenerObfsHttp.NewListener,
	consts.Otls:     listenerObfsTls.NewListener,
	consts.Pht:      listenerPht.NewListener,
	consts.Phts:     listenerPht.NewTLSListener,
	consts.Quic:     listenerQuic.NewListener,
	consts.Red:      listenerRedirectTcp.NewListener,
	consts.Redir:    listenerRedirectTcp.NewListener,
	consts.Redirect: listenerRedirectTcp.NewListener,
	consts.Redu:     listenerRedirectUdp.NewListener,
	consts.Rtcp:     listenerRtcp.NewListener,
	consts.Rudp:     listenerRudp.NewListener,
	consts.Ssh:      listenerSsh.NewListener,
	consts.Sshd:     listenerSshd.NewListener,
	consts.Tap:      listenerTap.NewListener,
	consts.Tcp:      listenerTcp.NewListener,
	consts.Tls:      listenerTls.NewListener,
	consts.Tun:      listenerTun.NewListener,
	consts.Udp:      listenerUdp.NewListener,
	consts.Ws:       listenerWs.NewListener,
	consts.Wss:      listenerWs.NewTLSListener,
}
