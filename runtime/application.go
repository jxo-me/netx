package runtime

import (
	"github.com/jxo-me/netx/core/admission"
	"github.com/jxo-me/netx/core/auth"
	"github.com/jxo-me/netx/core/bypass"
	"github.com/jxo-me/netx/core/chain"
	direct "github.com/jxo-me/netx/core/connector/direct"
	"github.com/jxo-me/netx/core/connector/forward"
	"github.com/jxo-me/netx/core/connector/http"
	"github.com/jxo-me/netx/core/connector/http2"
	"github.com/jxo-me/netx/core/connector/relay"
	"github.com/jxo-me/netx/core/connector/sni"
	v4 "github.com/jxo-me/netx/core/connector/socks/v4"
	v5 "github.com/jxo-me/netx/core/connector/socks/v5"
	"github.com/jxo-me/netx/core/connector/ss"
	ssu "github.com/jxo-me/netx/core/connector/ss/udp"
	"github.com/jxo-me/netx/core/connector/sshd"
	dialerDirect "github.com/jxo-me/netx/core/dialer/direct"
	"github.com/jxo-me/netx/core/dialer/dtls"
	"github.com/jxo-me/netx/core/dialer/ftcp"
	"github.com/jxo-me/netx/core/dialer/grpc"
	dialerHttp2 "github.com/jxo-me/netx/core/dialer/http2"
	"github.com/jxo-me/netx/core/dialer/http2/h2"
	"github.com/jxo-me/netx/core/dialer/http3"
	dialerIcmp "github.com/jxo-me/netx/core/dialer/icmp"
	"github.com/jxo-me/netx/core/dialer/kcp"
	"github.com/jxo-me/netx/core/dialer/mtls"
	"github.com/jxo-me/netx/core/dialer/mws"
	dialerObfsHttp "github.com/jxo-me/netx/core/dialer/obfs/http"
	dialerObfsTls "github.com/jxo-me/netx/core/dialer/obfs/tls"
	"github.com/jxo-me/netx/core/dialer/pht"
	dialerQuic "github.com/jxo-me/netx/core/dialer/quic"
	"github.com/jxo-me/netx/core/dialer/ssh"
	dialerSshd "github.com/jxo-me/netx/core/dialer/sshd"
	dialerTcp "github.com/jxo-me/netx/core/dialer/tcp"
	dialerTls "github.com/jxo-me/netx/core/dialer/tls"
	dialerUdp "github.com/jxo-me/netx/core/dialer/udp"
	"github.com/jxo-me/netx/core/dialer/ws"
	"github.com/jxo-me/netx/core/handler"
	"github.com/jxo-me/netx/core/handler/auto"
	"github.com/jxo-me/netx/core/handler/dns"
	"github.com/jxo-me/netx/core/handler/forward/local"
	"github.com/jxo-me/netx/core/handler/forward/remote"
	handlerHttp "github.com/jxo-me/netx/core/handler/http"
	handlerHttp2 "github.com/jxo-me/netx/core/handler/http2"
	handlerHttp3 "github.com/jxo-me/netx/core/handler/http3"
	redirect "github.com/jxo-me/netx/core/handler/redirect/tcp"
	redirectUdp "github.com/jxo-me/netx/core/handler/redirect/udp"
	handlerRelay "github.com/jxo-me/netx/core/handler/relay"
	handlerSni "github.com/jxo-me/netx/core/handler/sni"
	handlerSocksV4 "github.com/jxo-me/netx/core/handler/socks/v4"
	handlerSocksV5 "github.com/jxo-me/netx/core/handler/socks/v5"
	handlerSs "github.com/jxo-me/netx/core/handler/ss"
	handlerSsUdp "github.com/jxo-me/netx/core/handler/ss/udp"
	handlerSshd "github.com/jxo-me/netx/core/handler/sshd"
	"github.com/jxo-me/netx/core/handler/tap"
	"github.com/jxo-me/netx/core/handler/tun"
	"github.com/jxo-me/netx/core/hosts"
	"github.com/jxo-me/netx/core/ingress"
	"github.com/jxo-me/netx/core/limiter/conn"
	"github.com/jxo-me/netx/core/limiter/rate"
	"github.com/jxo-me/netx/core/limiter/traffic"
	listenerDns "github.com/jxo-me/netx/core/listener/dns"
	listenerDtls "github.com/jxo-me/netx/core/listener/dtls"
	listenerFtcp "github.com/jxo-me/netx/core/listener/ftcp"
	listenerGrpc "github.com/jxo-me/netx/core/listener/grpc"
	listenerHttp2 "github.com/jxo-me/netx/core/listener/http2"
	listenerHttpH2 "github.com/jxo-me/netx/core/listener/http2/h2"
	listenerHttp3 "github.com/jxo-me/netx/core/listener/http3"
	listenerHttpH3 "github.com/jxo-me/netx/core/listener/http3/h3"
	listenerIcmp "github.com/jxo-me/netx/core/listener/icmp"
	listenerKcp "github.com/jxo-me/netx/core/listener/kcp"
	listenerMtls "github.com/jxo-me/netx/core/listener/mtls"
	listenerMws "github.com/jxo-me/netx/core/listener/mws"
	listenerObfsHttp "github.com/jxo-me/netx/core/listener/obfs/http"
	listenerObfsTls "github.com/jxo-me/netx/core/listener/obfs/tls"
	listenerPht "github.com/jxo-me/netx/core/listener/pht"
	listenerQuic "github.com/jxo-me/netx/core/listener/quic"
	listenerRedirectTcp "github.com/jxo-me/netx/core/listener/redirect/tcp"
	listenerRedirectUdp "github.com/jxo-me/netx/core/listener/redirect/udp"
	listenerRtcp "github.com/jxo-me/netx/core/listener/rtcp"
	listenerRudp "github.com/jxo-me/netx/core/listener/rudp"
	listenerSsh "github.com/jxo-me/netx/core/listener/ssh"
	listenerSshd "github.com/jxo-me/netx/core/listener/sshd"
	listenerTap "github.com/jxo-me/netx/core/listener/tap"
	listenerTcp "github.com/jxo-me/netx/core/listener/tcp"
	listenerTls "github.com/jxo-me/netx/core/listener/tls"
	listenerTun "github.com/jxo-me/netx/core/listener/tun"
	listenerUdp "github.com/jxo-me/netx/core/listener/udp"
	listenerWs "github.com/jxo-me/netx/core/listener/ws"
	"github.com/jxo-me/netx/core/recorder"
	"github.com/jxo-me/netx/core/registry"
	"github.com/jxo-me/netx/core/resolver"
	"github.com/jxo-me/netx/core/service"
)

type Application struct {
	admissionReg      registry.IRegistry[admission.IAdmission]
	autherReg         registry.IRegistry[auth.IAuthenticator]
	bypassReg         registry.IRegistry[bypass.IBypass]
	chainReg          registry.IRegistry[chain.IChainer]
	connectorReg      registry.IRegistry[registry.NewConnector]
	connLimiterReg    registry.IRegistry[conn.IConnLimiter]
	dialerReg         registry.IRegistry[registry.NewDialer]
	handlerReg        registry.IRegistry[registry.NewHandler]
	hopReg            registry.IRegistry[chain.IHop]
	hostsReg          registry.IRegistry[hosts.IHostMapper]
	ingressReg        registry.IRegistry[ingress.IIngress]
	listenerReg       registry.IRegistry[registry.NewListener]
	rateLimiterReg    registry.IRegistry[rate.IRateLimiter]
	recorderReg       registry.IRegistry[recorder.IRecorder]
	resolverReg       registry.IRegistry[resolver.IResolver]
	serviceReg        registry.IRegistry[service.IService]
	trafficLimiterReg registry.IRegistry[traffic.ITrafficLimiter]
}

func NewConfig() *Application {
	app := Application{
		admissionReg:      new(registry.AdmissionRegistry),
		autherReg:         new(registry.AutherRegistry),
		bypassReg:         new(registry.BypassRegistry),
		chainReg:          new(registry.ChainRegistry),
		connectorReg:      new(registry.ConnectorRegistry),
		connLimiterReg:    new(registry.ConnLimiterRegistry),
		dialerReg:         new(registry.DialerRegistry),
		handlerReg:        new(registry.HandlerRegistry),
		hopReg:            new(registry.HopRegistry),
		hostsReg:          new(registry.HostsRegistry),
		ingressReg:        new(registry.IngressRegistry),
		listenerReg:       new(registry.ListenerRegistry),
		rateLimiterReg:    new(registry.RateLimiterRegistry),
		recorderReg:       new(registry.RecorderRegistry),
		resolverReg:       new(registry.ResolverRegistry),
		serviceReg:        new(registry.ServiceRegistry),
		trafficLimiterReg: new(registry.TrafficLimiterRegistry),
	}

	// Register connectors
	app.InitConnector()
	// Register dialers
	app.InitDialer()
	// Register handlers
	app.InitHandler()
	// Register listeners
	app.InitListener()

	return &app
}

func (a *Application) InitConnector() {
	// connector
	// direct
	a.ConnectorRegistry().Register("direct", direct.NewConnector)
	a.ConnectorRegistry().Register("virtual", direct.NewConnector)
	// forward
	a.ConnectorRegistry().Register("forward", forward.NewConnector)
	// http
	a.ConnectorRegistry().Register("http", http.NewConnector)
	// http2
	a.ConnectorRegistry().Register("http2", http2.NewConnector)
	// relay
	a.ConnectorRegistry().Register("relay", relay.NewConnector)
	// sni
	a.ConnectorRegistry().Register("sni", sni.NewConnector)
	// socks v4
	a.ConnectorRegistry().Register("socks4", v4.NewConnector)
	a.ConnectorRegistry().Register("socks4a", v4.NewConnector)
	// socks v5
	a.ConnectorRegistry().Register("socks5", v5.NewConnector)
	a.ConnectorRegistry().Register("socks", v5.NewConnector)
	// ss
	a.ConnectorRegistry().Register("ss", ss.NewConnector)
	a.ConnectorRegistry().Register("ssu", ssu.NewConnector)
	// sshd
	a.ConnectorRegistry().Register("sshd", sshd.NewConnector)

}

func (a *Application) InitDialer() {
	// dialer
	// direct
	a.DialerRegistry().Register("direct", dialerDirect.NewDialer)
	a.DialerRegistry().Register("virtual", dialerDirect.NewDialer)
	// dtls
	a.DialerRegistry().Register("dtls", dtls.NewDialer)
	// ftcp
	a.DialerRegistry().Register("ftcp", ftcp.NewDialer)
	// grpc
	a.DialerRegistry().Register("grpc", grpc.NewDialer)
	// http2
	a.DialerRegistry().Register("http2", dialerHttp2.NewDialer)
	a.DialerRegistry().Register("h2", h2.NewTLSDialer)
	a.DialerRegistry().Register("h2c", h2.NewDialer)
	// http3
	a.DialerRegistry().Register("http3", http3.NewDialer)
	a.DialerRegistry().Register("h3", http3.NewDialer)
	// icmp
	a.DialerRegistry().Register("icmp", dialerIcmp.NewDialer)
	// kcp
	a.DialerRegistry().Register("kcp", kcp.NewDialer)
	// mtls
	a.DialerRegistry().Register("mtls", mtls.NewDialer)
	// mws
	a.DialerRegistry().Register("mws", mws.NewDialer)
	a.DialerRegistry().Register("mwss", mws.NewTLSDialer)
	// obfs
	a.DialerRegistry().Register("ohttp", dialerObfsHttp.NewDialer)
	a.DialerRegistry().Register("otls", dialerObfsTls.NewDialer)
	// pht
	a.DialerRegistry().Register("pht", pht.NewDialer)
	a.DialerRegistry().Register("phts", pht.NewTLSDialer)
	// quic
	a.DialerRegistry().Register("quic", dialerQuic.NewDialer)
	// ssh
	a.DialerRegistry().Register("ssh", ssh.NewDialer)
	// sshd
	a.DialerRegistry().Register("sshd", dialerSshd.NewDialer)
	// tcp
	a.DialerRegistry().Register("tcp", dialerTcp.NewDialer)
	// tls
	a.DialerRegistry().Register("tls", dialerTls.NewDialer)
	// udp
	a.DialerRegistry().Register("udp", dialerUdp.NewDialer)
	// ws
	a.DialerRegistry().Register("ws", ws.NewDialer)
	a.DialerRegistry().Register("wss", ws.NewTLSDialer)

}

func (a *Application) InitHandler() {
	// handler
	// auto
	a.HandlerRegistry().Register("auto", func(opts ...handler.Option) handler.IHandler {
		options := handler.Options{}
		for _, opt := range opts {
			opt(&options)
		}
		h := auto.NewHandler(opts...)
		if f := a.HandlerRegistry().Get("http"); f != nil {
			v := append(opts,
				handler.LoggerOption(options.Logger.WithFields(map[string]any{"handler": "http"})))
			h.SetHttpHandler(f(v...))
		}
		if f := a.HandlerRegistry().Get("socks4"); f != nil {
			v := append(opts,
				handler.LoggerOption(options.Logger.WithFields(map[string]any{"handler": "socks4"})))
			h.SetSocks4Handler(f(v...))
		}
		if f := a.HandlerRegistry().Get("socks5"); f != nil {
			v := append(opts,
				handler.LoggerOption(options.Logger.WithFields(map[string]any{"handler": "socks5"})))
			h.SetSocks5Handler(f(v...))
		}
		return h
	})

	// dns
	a.HandlerRegistry().Register("dns", dns.NewHandler)
	// forward local
	a.HandlerRegistry().Register("tcp", local.NewHandler)
	a.HandlerRegistry().Register("udp", local.NewHandler)
	a.HandlerRegistry().Register("forward", local.NewHandler)
	// forward remote
	a.HandlerRegistry().Register("rtcp", remote.NewHandler)
	a.HandlerRegistry().Register("rudp", remote.NewHandler)
	// http
	a.HandlerRegistry().Register("http", handlerHttp.NewHandler)
	// http2
	a.HandlerRegistry().Register("http2", handlerHttp2.NewHandler)
	// http3
	a.HandlerRegistry().Register("http3", handlerHttp3.NewHandler)
	// redirect tcp
	a.HandlerRegistry().Register("red", redirect.NewHandler)
	a.HandlerRegistry().Register("redir", redirect.NewHandler)
	a.HandlerRegistry().Register("redirect", redirect.NewHandler)
	// redirect udp
	a.HandlerRegistry().Register("redu", redirectUdp.NewHandler)
	// relay
	a.HandlerRegistry().Register("relay", handlerRelay.NewHandler)
	// sni
	a.HandlerRegistry().Register("sni", handlerSni.NewHandler)
	// socks v4
	a.HandlerRegistry().Register("socks4", handlerSocksV4.NewHandler)
	a.HandlerRegistry().Register("socks4a", handlerSocksV4.NewHandler)
	// socks v5
	a.HandlerRegistry().Register("socks5", handlerSocksV5.NewHandler)
	a.HandlerRegistry().Register("socks", handlerSocksV5.NewHandler)
	// ss
	a.HandlerRegistry().Register("ss", handlerSs.NewHandler)
	a.HandlerRegistry().Register("ssu", handlerSsUdp.NewHandler)
	// sshd
	a.HandlerRegistry().Register("sshd", handlerSshd.NewHandler)
	// tap
	a.HandlerRegistry().Register("http", tap.NewHandler)
	// tun
	a.HandlerRegistry().Register("tun", tun.NewHandler)
}

func (a *Application) InitListener() {
	// listener
	// dns
	a.ListenerRegistry().Register("dns", listenerDns.NewListener)
	// dtls
	a.ListenerRegistry().Register("dtls", listenerDtls.NewListener)
	// ftcp
	a.ListenerRegistry().Register("ftcp", listenerFtcp.NewListener)
	// grpc
	a.ListenerRegistry().Register("grpc", listenerGrpc.NewListener)
	// http2
	a.ListenerRegistry().Register("http2", listenerHttp2.NewListener)
	// http2 h2
	a.ListenerRegistry().Register("h2c", listenerHttpH2.NewListener)
	a.ListenerRegistry().Register("h2", listenerHttpH2.NewTLSListener)
	// http3
	a.ListenerRegistry().Register("http3", listenerHttp3.NewListener)
	// http3 h3
	a.ListenerRegistry().Register("h3", listenerHttpH3.NewListener)
	// icmp
	a.ListenerRegistry().Register("icmp", listenerIcmp.NewListener)
	// kcp
	a.ListenerRegistry().Register("kcp", listenerKcp.NewListener)
	// mtls
	a.ListenerRegistry().Register("mtls", listenerMtls.NewListener)
	// mws
	a.ListenerRegistry().Register("mws", listenerMws.NewListener)
	a.ListenerRegistry().Register("mwss", listenerMws.NewTLSListener)
	// obfs http
	a.ListenerRegistry().Register("ohttp", listenerObfsHttp.NewListener)
	// obfs tls
	a.ListenerRegistry().Register("otls", listenerObfsTls.NewListener)
	// pht
	a.ListenerRegistry().Register("pht", listenerPht.NewListener)
	a.ListenerRegistry().Register("phts", listenerPht.NewTLSListener)
	// quic
	a.ListenerRegistry().Register("quic", listenerQuic.NewListener)
	// redirect tcp
	a.ListenerRegistry().Register("red", listenerRedirectTcp.NewListener)
	a.ListenerRegistry().Register("redir", listenerRedirectTcp.NewListener)
	a.ListenerRegistry().Register("redirect", listenerRedirectTcp.NewListener)
	// redirect udp
	a.ListenerRegistry().Register("redu", listenerRedirectUdp.NewListener)
	// rtcp
	a.ListenerRegistry().Register("rtcp", listenerRtcp.NewListener)
	// rudp
	a.ListenerRegistry().Register("rudp", listenerRudp.NewListener)
	// ssh
	a.ListenerRegistry().Register("ssh", listenerSsh.NewListener)
	// sshd
	a.ListenerRegistry().Register("sshd", listenerSshd.NewListener)
	// tap
	a.ListenerRegistry().Register("tap", listenerTap.NewListener)
	// tcp
	a.ListenerRegistry().Register("tcp", listenerTcp.NewListener)
	// tls
	a.ListenerRegistry().Register("tls", listenerTls.NewListener)
	// tun
	a.ListenerRegistry().Register("tun", listenerTun.NewListener)
	// udp
	a.ListenerRegistry().Register("udp", listenerUdp.NewListener)
	// ws
	a.ListenerRegistry().Register("ws", listenerWs.NewListener)
	a.ListenerRegistry().Register("wss", listenerWs.NewTLSListener)
}

func (a *Application) AdmissionRegistry() registry.IRegistry[admission.IAdmission] {
	return a.admissionReg
}

func (a *Application) AutherRegistry() registry.IRegistry[auth.IAuthenticator] {
	return a.autherReg
}

func (a *Application) BypassRegistry() registry.IRegistry[bypass.IBypass] {
	return a.bypassReg
}

func (a *Application) ChainRegistry() registry.IRegistry[chain.IChainer] {
	return a.chainReg
}

func (a *Application) ConnectorRegistry() registry.IRegistry[registry.NewConnector] {
	return a.connectorReg
}

func (a *Application) ConnLimiterRegistry() registry.IRegistry[conn.IConnLimiter] {
	return a.connLimiterReg
}

func (a *Application) DialerRegistry() registry.IRegistry[registry.NewDialer] {
	return a.dialerReg
}

func (a *Application) HandlerRegistry() registry.IRegistry[registry.NewHandler] {
	return a.handlerReg
}

func (a *Application) HopRegistry() registry.IRegistry[chain.IHop] {
	return a.hopReg
}

func (a *Application) HostsRegistry() registry.IRegistry[hosts.IHostMapper] {
	return a.hostsReg
}

func (a *Application) IngressRegistry() registry.IRegistry[ingress.IIngress] {
	return a.ingressReg
}

func (a *Application) ListenerRegistry() registry.IRegistry[registry.NewListener] {
	return a.listenerReg
}

func (a *Application) RateLimiterRegistry() registry.IRegistry[rate.IRateLimiter] {
	return a.rateLimiterReg
}

func (a *Application) RecorderRegistry() registry.IRegistry[recorder.IRecorder] {
	return a.recorderReg
}

func (a *Application) ResolverRegistry() registry.IRegistry[resolver.IResolver] {
	return a.resolverReg
}

func (a *Application) ServiceRegistry() registry.IRegistry[service.IService] {
	return a.serviceReg
}

func (a *Application) TrafficLimiterRegistry() registry.IRegistry[traffic.ITrafficLimiter] {
	return a.trafficLimiterReg
}
