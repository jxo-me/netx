package node

import (
	"fmt"
	"github.com/jxo-me/netx/x/app"
	"net"
	"strings"
	"time"

	"github.com/jxo-me/netx/core/bypass"
	"github.com/jxo-me/netx/core/chain"
	"github.com/jxo-me/netx/core/connector"
	"github.com/jxo-me/netx/core/dialer"
	"github.com/jxo-me/netx/core/logger"
	"github.com/jxo-me/netx/core/metadata"
	mdutil "github.com/jxo-me/netx/core/metadata/util"
	xauth "github.com/jxo-me/netx/x/auth"
	"github.com/jxo-me/netx/x/config"
	"github.com/jxo-me/netx/x/config/parsing"
	auth_parser "github.com/jxo-me/netx/x/config/parsing/auth"
	bypass_parser "github.com/jxo-me/netx/x/config/parsing/bypass"
	tls_util "github.com/jxo-me/netx/x/internal/util/tls"
	mdx "github.com/jxo-me/netx/x/metadata"
)

func ParseNode(hop string, cfg *config.NodeConfig, log logger.ILogger) (*chain.Node, error) {
	if cfg == nil {
		return nil, nil
	}

	if cfg.Connector == nil {
		cfg.Connector = &config.ConnectorConfig{
			Type: "http",
		}
	}

	if cfg.Dialer == nil {
		cfg.Dialer = &config.DialerConfig{
			Type: "tcp",
		}
	}

	nodeLogger := log.WithFields(map[string]any{
		"hop":       hop,
		"kind":      "node",
		"node":      cfg.Name,
		"connector": cfg.Connector.Type,
		"dialer":    cfg.Dialer.Type,
	})

	serverName, _, _ := net.SplitHostPort(cfg.Addr)

	tlsCfg := cfg.Connector.TLS
	if tlsCfg == nil {
		tlsCfg = &config.TLSConfig{}
	}
	if tlsCfg.ServerName == "" {
		tlsCfg.ServerName = serverName
	}
	tlsConfig, err := tls_util.LoadClientConfig(tlsCfg)
	if err != nil {
		nodeLogger.Error(err)
		return nil, err
	}

	var nm metadata.IMetaData
	if cfg.Metadata != nil {
		nm = mdx.NewMetadata(cfg.Metadata)
	}

	connectorLogger := nodeLogger.WithFields(map[string]any{
		"kind": "connector",
	})
	var cr connector.IConnector
	if rf := app.Runtime.ConnectorRegistry().Get(cfg.Connector.Type); rf != nil {
		cr = rf(
			connector.AuthOption(auth_parser.Info(cfg.Connector.Auth)),
			connector.TLSConfigOption(tlsConfig),
			connector.LoggerOption(connectorLogger),
		)
	} else {
		return nil, fmt.Errorf("unregistered connector: %s", cfg.Connector.Type)
	}

	if cfg.Connector.Metadata == nil {
		cfg.Connector.Metadata = make(map[string]any)
	}
	if err := cr.Init(mdx.NewMetadata(cfg.Connector.Metadata)); err != nil {
		connectorLogger.Error("init: ", err)
		return nil, err
	}

	tlsCfg = cfg.Dialer.TLS
	if tlsCfg == nil {
		tlsCfg = &config.TLSConfig{}
	}
	if tlsCfg.ServerName == "" {
		tlsCfg.ServerName = serverName
	}
	tlsConfig, err = tls_util.LoadClientConfig(tlsCfg)
	if err != nil {
		nodeLogger.Error(err)
		return nil, err
	}

	var ppv int
	if nm != nil {
		ppv = mdutil.GetInt(nm, parsing.MDKeyProxyProtocol)
	}

	dialerLogger := nodeLogger.WithFields(map[string]any{
		"kind": "dialer",
	})

	var d dialer.IDialer
	if rf := app.Runtime.DialerRegistry().Get(cfg.Dialer.Type); rf != nil {
		d = rf(
			dialer.AuthOption(auth_parser.Info(cfg.Dialer.Auth)),
			dialer.TLSConfigOption(tlsConfig),
			dialer.LoggerOption(dialerLogger),
			dialer.ProxyProtocolOption(ppv),
		)
	} else {
		return nil, fmt.Errorf("unregistered dialer: %s", cfg.Dialer.Type)
	}

	if cfg.Dialer.Metadata == nil {
		cfg.Dialer.Metadata = make(map[string]any)
	}
	if err := d.Init(mdx.NewMetadata(cfg.Dialer.Metadata)); err != nil {
		dialerLogger.Error("init: ", err)
		return nil, err
	}

	var sockOpts *chain.SockOpts
	if cfg.SockOpts != nil {
		sockOpts = &chain.SockOpts{
			Mark: cfg.SockOpts.Mark,
		}
	}

	tr := chain.NewTransport(d, cr,
		chain.AddrTransportOption(cfg.Addr),
		chain.InterfaceTransportOption(cfg.Interface),
		chain.SockOptsTransportOption(sockOpts),
		chain.TimeoutTransportOption(10*time.Second),
	)

	// convert *.example.com to .example.com
	// convert *example.com to example.com
	host := cfg.Host
	if strings.HasPrefix(host, "*") {
		host = host[1:]
		if !strings.HasPrefix(host, ".") {
			host = "." + host
		}
	}

	opts := []chain.NodeOption{
		chain.TransportNodeOption(tr),
		chain.BypassNodeOption(bypass.BypassGroup(bypass_parser.List(cfg.Bypass, cfg.Bypasses...)...)),
		chain.ResoloverNodeOption(app.Runtime.ResolverRegistry().Get(cfg.Resolver)),
		chain.HostMapperNodeOption(app.Runtime.HostsRegistry().Get(cfg.Hosts)),
		chain.MetadataNodeOption(nm),
		chain.HostNodeOption(host),
		chain.ProtocolNodeOption(cfg.Protocol),
		chain.PathNodeOption(cfg.Path),
		chain.NetworkNodeOption(cfg.Network),
	}
	if cfg.HTTP != nil {
		opts = append(opts, chain.HTTPNodeOption(&chain.HTTPNodeSettings{
			Host:   cfg.HTTP.Host,
			Header: cfg.HTTP.Header,
		}))
	}
	if cfg.TLS != nil {
		tlsCfg := &chain.TLSNodeSettings{
			ServerName: cfg.TLS.ServerName,
			Secure:     cfg.TLS.Secure,
		}
		if o := cfg.TLS.Options; o != nil {
			tlsCfg.Options.MinVersion = o.MinVersion
			tlsCfg.Options.MaxVersion = o.MaxVersion
			tlsCfg.Options.CipherSuites = o.CipherSuites
		}
		opts = append(opts, chain.TLSNodeOption(tlsCfg))
	}
	if cfg.Auth != nil {
		opts = append(opts, chain.AutherNodeOption(
			xauth.NewAuthenticator(
				xauth.AuthsOption(map[string]string{cfg.Auth.Username: cfg.Auth.Password}),
				xauth.LoggerOption(logger.Default().WithFields(map[string]any{
					"kind":     "node",
					"node":     cfg.Name,
					"addr":     cfg.Addr,
					"host":     cfg.Host,
					"protocol": cfg.Protocol,
				})),
			)))
	}
	return chain.NewNode(cfg.Name, cfg.Addr, opts...), nil
}