package boot

import (
	"github.com/jxo-me/netx/x/consts"
	listenerDns "github.com/jxo-me/netx/x/listener/dns"
	listenerDtls "github.com/jxo-me/netx/x/listener/dtls"
	listenerFtcp "github.com/jxo-me/netx/x/listener/ftcp"
	listenerGrpc "github.com/jxo-me/netx/x/listener/grpc"
	listenerHttp2 "github.com/jxo-me/netx/x/listener/http2"
	listenerHttpH2 "github.com/jxo-me/netx/x/listener/http2/h2"
	listenerHttp3 "github.com/jxo-me/netx/x/listener/http3"
	listenerHttpH3 "github.com/jxo-me/netx/x/listener/http3/h3"
	listenerHttpWt "github.com/jxo-me/netx/x/listener/http3/wt"
	listenerIcmp "github.com/jxo-me/netx/x/listener/icmp"
	listenerKcp "github.com/jxo-me/netx/x/listener/kcp"
	listenerMtcp "github.com/jxo-me/netx/x/listener/mtcp"
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
	listenerSerial "github.com/jxo-me/netx/x/listener/serial"
	listenerSsh "github.com/jxo-me/netx/x/listener/ssh"
	listenerSshd "github.com/jxo-me/netx/x/listener/sshd"
	listenerTap "github.com/jxo-me/netx/x/listener/tap"
	listenerTcp "github.com/jxo-me/netx/x/listener/tcp"
	listenerTls "github.com/jxo-me/netx/x/listener/tls"
	listenerTun "github.com/jxo-me/netx/x/listener/tun"
	listenerUdp "github.com/jxo-me/netx/x/listener/udp"
	listenerUnix "github.com/jxo-me/netx/x/listener/unix"
	listenerWs "github.com/jxo-me/netx/x/listener/ws"
	"github.com/jxo-me/netx/x/registry"
)

var Listeners = map[string]registry.NewListener{
	consts.Dns:      listenerDns.NewListener,
	consts.Dtls:     listenerDtls.NewListener,
	consts.Ftcp:     listenerFtcp.NewListener,
	consts.Grpc:     listenerGrpc.NewListener,
	consts.Http2:    listenerHttp2.NewListener,
	consts.H2c:      listenerHttpH2.NewListener,
	consts.H2:       listenerHttpH2.NewTLSListener,
	consts.Http3:    listenerHttp3.NewListener,
	consts.H3:       listenerHttpH3.NewListener,
	consts.Wt:       listenerHttpWt.NewListener,
	consts.Icmp:     listenerIcmp.NewListener,
	consts.Kcp:      listenerKcp.NewListener,
	consts.Mtcp:     listenerMtcp.NewListener,
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
	consts.Serial:   listenerSerial.NewListener,
	consts.Ssh:      listenerSsh.NewListener,
	consts.Sshd:     listenerSshd.NewListener,
	consts.Tap:      listenerTap.NewListener,
	consts.Tcp:      listenerTcp.NewListener,
	consts.Tls:      listenerTls.NewListener,
	consts.Tun:      listenerTun.NewListener,
	consts.Udp:      listenerUdp.NewListener,
	consts.Unix:     listenerUnix.NewListener,
	consts.Ws:       listenerWs.NewListener,
	consts.Wss:      listenerWs.NewTLSListener,
}
