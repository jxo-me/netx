package pht

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/jxo-me/netx/core/dialer"
	md "github.com/jxo-me/netx/core/metadata"
	pht_util "github.com/jxo-me/netx/x/internal/util/pht"
	"github.com/jxo-me/netx/x/registry"
)

func init() {
	registry.DialerRegistry().Register("pht", NewDialer)
	registry.DialerRegistry().Register("phts", NewTLSDialer)
}

type phtDialer struct {
	clients     map[string]*pht_util.Client
	clientMutex sync.Mutex
	tlsEnabled  bool
	md          metadata
	options     dialer.Options
}

func NewDialer(opts ...dialer.Option) dialer.IDialer {
	options := dialer.Options{}
	for _, opt := range opts {
		opt(&options)
	}

	return &phtDialer{
		clients: make(map[string]*pht_util.Client),
		options: options,
	}
}

func NewTLSDialer(opts ...dialer.Option) dialer.IDialer {
	options := dialer.Options{}
	for _, opt := range opts {
		opt(&options)
	}

	return &phtDialer{
		tlsEnabled: true,
		clients:    make(map[string]*pht_util.Client),
		options:    options,
	}
}

func (d *phtDialer) Init(md md.IMetaData) (err error) {
	if err = d.parseMetadata(md); err != nil {
		return
	}

	return nil
}

// Multiplex implements dialer.IMultiplexer interface.
// NOTE: PHT is not a real multiplexed tunnel.
func (d *phtDialer) Multiplex() bool {
	return true
}

func (d *phtDialer) Dial(ctx context.Context, addr string, opts ...dialer.DialOption) (net.Conn, error) {
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

		tr := &http.Transport{
			// Proxy: http.ProxyFromEnvironment,
			DialContext: func(ctx context.Context, network, adr string) (net.Conn, error) {
				return options.NetDialer.Dial(ctx, network, addr)
			},
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
		if d.tlsEnabled {
			tr.TLSClientConfig = d.options.TLSConfig
		}

		client = &pht_util.Client{
			Host: host,
			Client: &http.Client{
				// Timeout:   60 * time.Second,
				Transport: tr,
			},
			AuthorizePath: d.md.authorizePath,
			PushPath:      d.md.pushPath,
			PullPath:      d.md.pullPath,
			TLSEnabled:    d.tlsEnabled,
			Logger:        d.options.Logger,
		}
		d.clients[addr] = client
	}

	return client.Dial(ctx, addr)
}
