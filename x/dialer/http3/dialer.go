package http3

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"sync"

	"github.com/jxo-me/netx/core/dialer"
	md "github.com/jxo-me/netx/core/metadata"
	pht_util "github.com/jxo-me/netx/x/internal/util/pht"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

type http3Dialer struct {
	clients     map[string]*pht_util.Client
	clientMutex sync.Mutex
	md          metadata
	options     dialer.Options
}

func NewDialer(opts ...dialer.Option) dialer.IDialer {
	options := dialer.Options{}
	for _, opt := range opts {
		opt(&options)
	}

	return &http3Dialer{
		clients: make(map[string]*pht_util.Client),
		options: options,
	}
}

func (d *http3Dialer) Init(md md.IMetaData) (err error) {
	if err = d.parseMetadata(md); err != nil {
		return
	}

	return nil
}

func (d *http3Dialer) Dial(ctx context.Context, addr string, opts ...dialer.DialOption) (net.Conn, error) {
	d.clientMutex.Lock()
	defer d.clientMutex.Unlock()

	client, ok := d.clients[addr]
	if !ok {
		var options dialer.DialOptions
		for _, opt := range opts {
			opt(&options)
		}

		host := d.md.host
		if host == "" {
			host = options.Host
		}
		if h, _, _ := net.SplitHostPort(host); h != "" {
			host = h
		}

		client = &pht_util.Client{
			Host: host,
			Client: &http.Client{
				// Timeout:   60 * time.Second,
				Transport: &http3.RoundTripper{
					TLSClientConfig: d.options.TLSConfig,
					Dial: func(ctx context.Context, adr string, tlsCfg *tls.Config, cfg *quic.Config) (quic.EarlyConnection, error) {
						// d.options.Logger.Infof("dial: %s/%s, %s", addr, network, host)
						udpAddr, err := net.ResolveUDPAddr("udp", addr)
						if err != nil {
							return nil, err
						}

						udpConn, err := options.Dialer.Dial(ctx, "udp", "")
						if err != nil {
							return nil, err
						}

						return quic.DialEarly(context.Background(), udpConn.(net.PacketConn), udpAddr, tlsCfg, cfg)
					},
					QUICConfig: &quic.Config{
						KeepAlivePeriod:      d.md.keepAlivePeriod,
						HandshakeIdleTimeout: d.md.handshakeTimeout,
						MaxIdleTimeout:       d.md.maxIdleTimeout,
						Versions: []quic.Version{
							quic.Version1,
							quic.Version2,
						},
						MaxIncomingStreams: int64(d.md.maxStreams),
					},
				},
			},
			AuthorizePath: d.md.authorizePath,
			PushPath:      d.md.pushPath,
			PullPath:      d.md.pullPath,
			TLSEnabled:    true,
			Logger:        d.options.Logger,
		}

		d.clients[addr] = client
	}

	return client.Dial(ctx, addr)
}

// Multiplex implements dialer.IMultiplexer interface.
func (d *http3Dialer) Multiplex() bool {
	return true
}
