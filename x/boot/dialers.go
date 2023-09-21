package boot

import (
	"github.com/jxo-me/netx/core/dialer"
	"github.com/jxo-me/netx/x/consts"
	dialerDirect "github.com/jxo-me/netx/x/dialer/direct"
	"github.com/jxo-me/netx/x/dialer/dtls"
	"github.com/jxo-me/netx/x/dialer/ftcp"
	"github.com/jxo-me/netx/x/dialer/grpc"
	dialerHttp2 "github.com/jxo-me/netx/x/dialer/http2"
	"github.com/jxo-me/netx/x/dialer/http2/h2"
	"github.com/jxo-me/netx/x/dialer/http3"
	dialerIcmp "github.com/jxo-me/netx/x/dialer/icmp"
	"github.com/jxo-me/netx/x/dialer/kcp"
	"github.com/jxo-me/netx/x/dialer/mtls"
	"github.com/jxo-me/netx/x/dialer/mws"
	dialerObfsHttp "github.com/jxo-me/netx/x/dialer/obfs/http"
	dialerObfsTls "github.com/jxo-me/netx/x/dialer/obfs/tls"
	"github.com/jxo-me/netx/x/dialer/pht"
	dialerQuic "github.com/jxo-me/netx/x/dialer/quic"
	dialerSerial "github.com/jxo-me/netx/x/dialer/serial"
	"github.com/jxo-me/netx/x/dialer/ssh"
	dialerSshd "github.com/jxo-me/netx/x/dialer/sshd"
	dialerTcp "github.com/jxo-me/netx/x/dialer/tcp"
	dialerTls "github.com/jxo-me/netx/x/dialer/tls"
	dialerUdp "github.com/jxo-me/netx/x/dialer/udp"
	dialerUnix "github.com/jxo-me/netx/x/dialer/unix"
	"github.com/jxo-me/netx/x/dialer/ws"
)

var Dialers = map[string]dialer.NewDialer{
	consts.Direct:  dialerDirect.NewDialer,
	consts.Virtual: dialerDirect.NewDialer,
	consts.Dtls:    dtls.NewDialer,
	consts.Ftcp:    ftcp.NewDialer,
	consts.Grpc:    grpc.NewDialer,
	consts.Http2:   dialerHttp2.NewDialer,
	consts.H2:      h2.NewTLSDialer,
	consts.H2c:     h2.NewDialer,
	consts.Http3:   http3.NewDialer,
	consts.H3:      http3.NewDialer,
	consts.Icmp:    dialerIcmp.NewDialer,
	consts.Kcp:     kcp.NewDialer,
	consts.Mtls:    mtls.NewDialer,
	consts.Mws:     mws.NewDialer,
	consts.Mwss:    mws.NewTLSDialer,
	consts.Ohttp:   dialerObfsHttp.NewDialer,
	consts.Otls:    dialerObfsTls.NewDialer,
	consts.Pht:     pht.NewDialer,
	consts.Phts:    pht.NewTLSDialer,
	consts.Quic:    dialerQuic.NewDialer,
	consts.Serial:  dialerSerial.NewDialer,
	consts.Ssh:     ssh.NewDialer,
	consts.Sshd:    dialerSshd.NewDialer,
	consts.Tcp:     dialerTcp.NewDialer,
	consts.Tls:     dialerTls.NewDialer,
	consts.Udp:     dialerUdp.NewDialer,
	consts.Unix:    dialerUnix.NewDialer,
	consts.Ws:      ws.NewDialer,
	consts.Wss:     ws.NewTLSDialer,
}
