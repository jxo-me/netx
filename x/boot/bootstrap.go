package boot

import (
	"fmt"
	"github.com/jxo-me/netx/core/handler"
	"github.com/jxo-me/netx/x/app"
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
	"github.com/jxo-me/netx/x/dialer/ssh"
	dialerSshd "github.com/jxo-me/netx/x/dialer/sshd"
	dialerTcp "github.com/jxo-me/netx/x/dialer/tcp"
	dialerTls "github.com/jxo-me/netx/x/dialer/tls"
	dialerUdp "github.com/jxo-me/netx/x/dialer/udp"
	"github.com/jxo-me/netx/x/dialer/ws"
	"github.com/jxo-me/netx/x/handler/auto"
	"github.com/jxo-me/netx/x/handler/dns"
	"github.com/jxo-me/netx/x/handler/forward/local"
	"github.com/jxo-me/netx/x/handler/forward/remote"
	handlerHttp "github.com/jxo-me/netx/x/handler/http"
	handlerHttp2 "github.com/jxo-me/netx/x/handler/http2"
	handlerHttp3 "github.com/jxo-me/netx/x/handler/http3"
	redirect "github.com/jxo-me/netx/x/handler/redirect/tcp"
	redirectUdp "github.com/jxo-me/netx/x/handler/redirect/udp"
	handlerRelay "github.com/jxo-me/netx/x/handler/relay"
	handlerSni "github.com/jxo-me/netx/x/handler/sni"
	handlerSocksV4 "github.com/jxo-me/netx/x/handler/socks/v4"
	handlerSocksV5 "github.com/jxo-me/netx/x/handler/socks/v5"
	handlerSs "github.com/jxo-me/netx/x/handler/ss"
	handlerSsUdp "github.com/jxo-me/netx/x/handler/ss/udp"
	handlerSshd "github.com/jxo-me/netx/x/handler/sshd"
	"github.com/jxo-me/netx/x/handler/tap"
	"github.com/jxo-me/netx/x/handler/tun"
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
	"github.com/jxo-me/netx/x/registry"
)

var (
	insBoot = Boot{}
)

func Boots(a app.IRuntime) *Boot {
	insBoot.App = a
	insBoot.Connectors = map[string]registry.NewConnector{
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
	insBoot.Dialers = map[string]registry.NewDialer{
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
		consts.Ssh:     ssh.NewDialer,
		consts.Sshd:    dialerSshd.NewDialer,
		consts.Tcp:     dialerTcp.NewDialer,
		consts.Tls:     dialerTls.NewDialer,
		consts.Udp:     dialerUdp.NewDialer,
		consts.Ws:      ws.NewDialer,
		consts.Wss:     ws.NewTLSDialer,
	}
	insBoot.Handlers = map[string]registry.NewHandler{
		consts.Dns:      dns.NewHandler,
		consts.Tcp:      local.NewHandler,
		consts.Udp:      local.NewHandler,
		consts.Forward:  local.NewHandler,
		consts.Rtcp:     remote.NewHandler,
		consts.Rudp:     remote.NewHandler,
		consts.Http:     handlerHttp.NewHandler,
		consts.Http2:    handlerHttp2.NewHandler,
		consts.Http3:    handlerHttp3.NewHandler,
		consts.Red:      redirect.NewHandler,
		consts.Redir:    redirect.NewHandler,
		consts.Redirect: redirect.NewHandler,
		consts.Redu:     redirectUdp.NewHandler,
		consts.Relay:    handlerRelay.NewHandler,
		consts.Sni:      handlerSni.NewHandler,
		consts.Socks4:   handlerSocksV4.NewHandler,
		consts.Socks4a:  handlerSocksV4.NewHandler,
		consts.Socks5:   handlerSocksV5.NewHandler,
		consts.Socks:    handlerSocksV5.NewHandler,
		consts.Ss:       handlerSs.NewHandler,
		consts.Ssu:      handlerSsUdp.NewHandler,
		consts.Sshd:     handlerSshd.NewHandler,
		consts.Tap:      tap.NewHandler,
		consts.Tun:      tun.NewHandler,
	}
	insBoot.Listeners = map[string]registry.NewListener{
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
	// Register connectors
	err := insBoot.InitConnector()
	if err != nil {
		panic(fmt.Sprintf("InitConnector error: %s", err.Error()))
		return nil
	}
	// Register dialers
	err = insBoot.InitDialer()
	if err != nil {
		panic(fmt.Sprintf("InitDialer error: %s", err.Error()))
		return nil
	}
	// Register handlers
	err = insBoot.InitHandler()
	if err != nil {
		panic(fmt.Sprintf("InitHandler error: %s", err.Error()))
		return nil
	}
	// Register listeners
	err = insBoot.InitListener()
	if err != nil {
		panic(fmt.Sprintf("InitListener error: %s", err.Error()))
		return nil
	}

	return &insBoot
}

type Boot struct {
	App        app.IRuntime
	Connectors map[string]registry.NewConnector
	Dialers    map[string]registry.NewDialer
	Handlers   map[string]registry.NewHandler
	Listeners  map[string]registry.NewListener
}

func (b *Boot) InitConnector() (err error) {
	// connector
	for name, connector := range b.Connectors {
		//fmt.Println("Register Connector type:", name)
		err = b.App.ConnectorRegistry().Register(name, connector)
		if err != nil {
			return
		}
	}
	return err
}

func (b *Boot) InitDialer() (err error) {
	// dialer
	for name, dialer := range b.Dialers {
		//fmt.Println("Register Dialer type:", name)
		err = b.App.DialerRegistry().Register(name, dialer)
		if err != nil {
			return err
		}
	}
	return err
}

func (b *Boot) InitHandler() (err error) {
	// handler
	for name, handle := range b.Handlers {
		//fmt.Println("Register Handler type:", name)
		if name == consts.Auto {
			err = b.App.HandlerRegistry().Register(consts.Auto, func(opts ...handler.Option) handler.IHandler {
				options := handler.Options{}
				for _, opt := range opts {
					opt(&options)
				}
				h := auto.NewHandler(opts...)
				if f := b.App.HandlerRegistry().Get(consts.Http); f != nil {
					v := append(opts,
						handler.LoggerOption(options.Logger.WithFields(map[string]any{"handler": consts.Http})))
					h.SetHttpHandler(f(v...))
				}
				if f := b.App.HandlerRegistry().Get(consts.Socks4); f != nil {
					v := append(opts,
						handler.LoggerOption(options.Logger.WithFields(map[string]any{"handler": consts.Socks4})))
					h.SetSocks4Handler(f(v...))
				}
				if f := b.App.HandlerRegistry().Get(consts.Socks5); f != nil {
					v := append(opts,
						handler.LoggerOption(options.Logger.WithFields(map[string]any{"handler": consts.Socks5})))
					h.SetSocks5Handler(f(v...))
				}
				return h
			})
			if err != nil {
				return err
			}
		} else {
			err = b.App.HandlerRegistry().Register(name, handle)
			if err != nil {
				return err
			}
		}
	}
	return err
}

func (b *Boot) InitListener() (err error) {
	// listener
	for name, listener := range b.Listeners {
		//fmt.Println("Register Listener type:", name)
		err = b.App.ListenerRegistry().Register(name, listener)
		if err != nil {
			return err
		}
	}
	return err
}
