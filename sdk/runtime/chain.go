package runtime

import (
	"fmt"
	"net"
	"strings"
	"time"

	auther "github.com/jxo-me/netx/core/auth"
	"github.com/jxo-me/netx/core/bypass"
	"github.com/jxo-me/netx/core/chain"
	xchain "github.com/jxo-me/netx/core/chain"
	"github.com/jxo-me/netx/core/connector"
	"github.com/jxo-me/netx/core/dialer"
	"github.com/jxo-me/netx/core/logger"
	"github.com/jxo-me/netx/core/metadata"
	mdx "github.com/jxo-me/netx/core/metadata"
	mdutil "github.com/jxo-me/netx/core/metadata/util"
	tls_util "github.com/jxo-me/netx/internal/util/tls"
	"github.com/jxo-me/netx/x/config"
)

func (a *Application) ParseChain(cfg *config.ChainConfig) (chain.IChainer, error) {
	if cfg == nil {
		return nil, nil
	}

	chainLogger := logger.Default().WithFields(map[string]any{
		"kind":  "chain",
		"chain": cfg.Name,
	})

	var md metadata.IMetaData
	if cfg.Metadata != nil {
		md = mdx.NewMetadata(cfg.Metadata)
	}

	c := xchain.NewChain(cfg.Name,
		xchain.MetadataChainOption(md),
		xchain.LoggerChainOption(chainLogger),
	)

	for _, ch := range cfg.Hops {
		var hop chain.IHop
		var err error

		if len(ch.Nodes) > 0 {
			if hop, err = a.ParseHop(ch); err != nil {
				return nil, err
			}
		} else {
			hop = a.HopRegistry().Get(ch.Name)
		}
		if hop != nil {
			c.AddHop(hop)
		}
	}

	return c, nil
}

func (a *Application) ParseHop(cfg *config.HopConfig) (chain.IHop, error) {
	if cfg == nil {
		return nil, nil
	}

	hopLogger := logger.Default().WithFields(map[string]any{
		"kind": "hop",
		"hop":  cfg.Name,
	})

	var nodes []*chain.Node
	for _, v := range cfg.Nodes {
		if v == nil {
			continue
		}

		if v.Connector == nil {
			v.Connector = &config.ConnectorConfig{
				Type: "http",
			}
		}

		if v.Dialer == nil {
			v.Dialer = &config.DialerConfig{
				Type: "tcp",
			}
		}

		nodeLogger := hopLogger.WithFields(map[string]any{
			"kind":      "node",
			"node":      v.Name,
			"connector": v.Connector.Type,
			"dialer":    v.Dialer.Type,
		})

		serverName, _, _ := net.SplitHostPort(v.Addr)

		tlsCfg := v.Connector.TLS
		if tlsCfg == nil {
			tlsCfg = &config.TLSConfig{}
		}
		if tlsCfg.ServerName == "" {
			tlsCfg.ServerName = serverName
		}
		tlsConfig, err := tls_util.LoadClientConfig(
			tlsCfg.CertFile, tlsCfg.KeyFile, tlsCfg.CAFile,
			tlsCfg.Secure, tlsCfg.ServerName)
		if err != nil {
			hopLogger.Error(err)
			return nil, err
		}

		var nm metadata.IMetaData
		if v.Metadata != nil {
			nm = mdx.NewMetadata(v.Metadata)
		}

		connectorLogger := nodeLogger.WithFields(map[string]any{
			"kind": "connector",
		})
		var cr connector.IConnector
		if rf := a.ConnectorRegistry().Get(v.Connector.Type); rf != nil {
			cr = rf(
				connector.AuthOption(a.parseAuth(v.Connector.Auth)),
				connector.TLSConfigOption(tlsConfig),
				connector.LoggerOption(connectorLogger),
			)
		} else {
			return nil, fmt.Errorf("unregistered connector: %s", v.Connector.Type)
		}

		if v.Connector.Metadata == nil {
			v.Connector.Metadata = make(map[string]any)
		}
		if err := cr.Init(mdx.NewMetadata(v.Connector.Metadata)); err != nil {
			connectorLogger.Error("init: ", err)
			return nil, err
		}

		tlsCfg = v.Dialer.TLS
		if tlsCfg == nil {
			tlsCfg = &config.TLSConfig{}
		}
		if tlsCfg.ServerName == "" {
			tlsCfg.ServerName = serverName
		}
		tlsConfig, err = tls_util.LoadClientConfig(
			tlsCfg.CertFile, tlsCfg.KeyFile, tlsCfg.CAFile,
			tlsCfg.Secure, tlsCfg.ServerName)
		if err != nil {
			hopLogger.Error(err)
			return nil, err
		}

		var ppv int
		if nm != nil {
			ppv = mdutil.GetInt(nm, mdKeyProxyProtocol)
		}

		dialerLogger := nodeLogger.WithFields(map[string]any{
			"kind": "dialer",
		})

		var d dialer.IDialer
		if rf := a.DialerRegistry().Get(v.Dialer.Type); rf != nil {
			d = rf(
				dialer.AuthOption(a.parseAuth(v.Dialer.Auth)),
				dialer.TLSConfigOption(tlsConfig),
				dialer.LoggerOption(dialerLogger),
				dialer.ProxyProtocolOption(ppv),
			)
		} else {
			return nil, fmt.Errorf("unregistered dialer: %s", v.Dialer.Type)
		}

		if v.Dialer.Metadata == nil {
			v.Dialer.Metadata = make(map[string]any)
		}
		if err := d.Init(mdx.NewMetadata(v.Dialer.Metadata)); err != nil {
			dialerLogger.Error("init: ", err)
			return nil, err
		}

		if v.Resolver == "" {
			v.Resolver = cfg.Resolver
		}
		if v.Hosts == "" {
			v.Hosts = cfg.Hosts
		}
		if v.Interface == "" {
			v.Interface = cfg.Interface
		}
		if v.SockOpts == nil {
			v.SockOpts = cfg.SockOpts
		}

		var sockOpts *chain.SockOpts
		if v.SockOpts != nil {
			sockOpts = &chain.SockOpts{
				Mark: v.SockOpts.Mark,
			}
		}

		tr := chain.NewTransport(d, cr,
			chain.AddrTransportOption(v.Addr),
			chain.InterfaceTransportOption(v.Interface),
			chain.SockOptsTransportOption(sockOpts),
			chain.TimeoutTransportOption(10*time.Second),
		)

		// convert *.example.com to .example.com
		// convert *example.com to example.com
		host := v.Host
		if strings.HasPrefix(host, "*") {
			host = host[1:]
			if !strings.HasPrefix(host, ".") {
				host = "." + host
			}
		}

		opts := []chain.NodeOption{
			chain.TransportNodeOption(tr),
			chain.BypassNodeOption(bypass.BypassGroup(a.bypassList(v.Bypass, v.Bypasses...)...)),
			chain.ResoloverNodeOption(a.ResolverRegistry().Get(v.Resolver)),
			chain.HostMapperNodeOption(a.HostsRegistry().Get(v.Hosts)),
			chain.MetadataNodeOption(nm),
			chain.HostNodeOption(host),
			chain.ProtocolNodeOption(v.Protocol),
		}
		if v.HTTP != nil {
			opts = append(opts, chain.HTTPNodeOption(&chain.HTTPNodeSettings{
				Host:   v.HTTP.Host,
				Header: v.HTTP.Header,
			}))
		}
		if v.TLS != nil {
			opts = append(opts, chain.TLSNodeOption(&chain.TLSNodeSettings{
				ServerName: v.TLS.ServerName,
				Secure:     v.TLS.Secure,
			}))
		}
		if v.Auth != nil {
			opts = append(opts, chain.AutherNodeOption(
				auther.NewAuthenticator(
					auther.AuthsOption(map[string]string{v.Auth.Username: v.Auth.Password}),
					auther.LoggerOption(logger.Default().WithFields(map[string]any{
						"kind":     "node",
						"node":     v.Name,
						"addr":     v.Addr,
						"host":     v.Host,
						"protocol": v.Protocol,
					})),
				)))
		}
		node := chain.NewNode(v.Name, v.Addr, opts...)
		nodes = append(nodes, node)
	}

	sel := a.parseNodeSelector(cfg.Selector)
	if sel == nil {
		sel = a.defaultNodeSelector()
	}
	return xchain.NewChainHop(nodes,
		xchain.SelectorHopOption(sel),
		xchain.BypassHopOption(bypass.BypassGroup(a.bypassList(cfg.Bypass, cfg.Bypasses...)...)),
		xchain.LoggerHopOption(hopLogger),
	), nil
}
