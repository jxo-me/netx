package runtime

import (
	"fmt"
	"strings"

	"github.com/jxo-me/netx/config"
	"github.com/jxo-me/netx/core/admission"
	"github.com/jxo-me/netx/core/auth"
	"github.com/jxo-me/netx/core/bypass"
	"github.com/jxo-me/netx/core/chain"
	xchain "github.com/jxo-me/netx/core/chain"
	"github.com/jxo-me/netx/core/handler"
	"github.com/jxo-me/netx/core/listener"
	"github.com/jxo-me/netx/core/logger"
	"github.com/jxo-me/netx/core/metadata"
	mdutil "github.com/jxo-me/netx/core/metadata/util"
	"github.com/jxo-me/netx/core/recorder"
	"github.com/jxo-me/netx/core/selector"
	"github.com/jxo-me/netx/core/service"
	xservice "github.com/jxo-me/netx/core/service"
	tls_util "github.com/jxo-me/netx/internal/util/tls"
)

func (a *Application) ParseService(cfg *config.ServiceConfig) (service.IService, error) {
	if cfg.Listener == nil {
		cfg.Listener = &config.ListenerConfig{
			Type: "tcp",
		}
	}
	if cfg.Handler == nil {
		cfg.Handler = &config.HandlerConfig{
			Type: "auto",
		}
	}
	serviceLogger := logger.Default().WithFields(map[string]any{
		"kind":     "service",
		"service":  cfg.Name,
		"listener": cfg.Listener.Type,
		"handler":  cfg.Handler.Type,
	})

	listenerLogger := serviceLogger.WithFields(map[string]any{
		"kind": "listener",
	})

	tlsCfg := cfg.Listener.TLS
	if tlsCfg == nil {
		tlsCfg = &config.TLSConfig{}
	}
	tlsConfig, err := tls_util.LoadServerConfig(
		tlsCfg.CertFile, tlsCfg.KeyFile, tlsCfg.CAFile)
	if err != nil {
		listenerLogger.Error(err)
		return nil, err
	}
	if tlsConfig == nil {
		tlsConfig = defaultTLSConfig.Clone()
	}

	authers := a.autherList(cfg.Listener.Auther, cfg.Listener.Authers...)
	if len(authers) == 0 {
		if auther := a.ParseAutherFromAuth(cfg.Listener.Auth); auther != nil {
			authers = append(authers, auther)
		}
	}
	var auther auth.IAuthenticator
	if len(authers) > 0 {
		auther = auth.AuthenticatorGroup(authers...)
	}

	admissions := a.admissionList(cfg.Admission, cfg.Admissions...)

	var sockOpts *chain.SockOpts
	if cfg.SockOpts != nil {
		sockOpts = &chain.SockOpts{
			Mark: cfg.SockOpts.Mark,
		}
	}

	var ppv int
	ifce := cfg.Interface
	var preUp, preDown, postUp, postDown []string
	var ignoreChain bool
	if cfg.Metadata != nil {
		md := metadata.NewMetadata(cfg.Metadata)
		ppv = mdutil.GetInt(md, mdKeyProxyProtocol)
		if v := mdutil.GetString(md, mdKeyInterface); v != "" {
			ifce = v
		}
		if v := mdutil.GetInt(md, mdKeySoMark); v > 0 {
			sockOpts = &chain.SockOpts{
				Mark: v,
			}
		}
		preUp = mdutil.GetStrings(md, mdKeyPreUp)
		preDown = mdutil.GetStrings(md, mdKeyPreDown)
		postUp = mdutil.GetStrings(md, mdKeyPostUp)
		postDown = mdutil.GetStrings(md, mdKeyPostDown)
		ignoreChain = mdutil.GetBool(md, mdKeyIgnoreChain)
	}

	listenOpts := []listener.Option{
		listener.AddrOption(cfg.Addr),
		listener.AutherOption(auther),
		listener.AuthOption(a.parseAuth(cfg.Listener.Auth)),
		listener.TLSConfigOption(tlsConfig),
		listener.AdmissionOption(admission.AdmissionGroup(admissions...)),
		listener.TrafficLimiterOption(a.TrafficLimiterRegistry().Get(cfg.Limiter)),
		listener.ConnLimiterOption(a.ConnLimiterRegistry().Get(cfg.CLimiter)),
		listener.LoggerOption(listenerLogger),
		listener.ServiceOption(cfg.Name),
		listener.ProxyProtocolOption(ppv),
	}
	if !ignoreChain {
		listenOpts = append(listenOpts,
			listener.ChainOption(a.chainGroup(cfg.Listener.Chain, cfg.Listener.ChainGroup)),
		)
	}

	var ln listener.IListener
	if rf := a.ListenerRegistry().Get(cfg.Listener.Type); rf != nil {
		ln = rf(listenOpts...)
	} else {
		return nil, fmt.Errorf("unregistered listener: %s", cfg.Listener.Type)
	}

	if cfg.Listener.Metadata == nil {
		cfg.Listener.Metadata = make(map[string]any)
	}
	listenerLogger.Debugf("metadata: %v", cfg.Listener.Metadata)
	if err := ln.Init(metadata.NewMetadata(cfg.Listener.Metadata)); err != nil {
		listenerLogger.Error("init: ", err)
		return nil, err
	}

	handlerLogger := serviceLogger.WithFields(map[string]any{
		"kind": "handler",
	})

	tlsCfg = cfg.Handler.TLS
	if tlsCfg == nil {
		tlsCfg = &config.TLSConfig{}
	}
	tlsConfig, err = tls_util.LoadServerConfig(
		tlsCfg.CertFile, tlsCfg.KeyFile, tlsCfg.CAFile)
	if err != nil {
		handlerLogger.Error(err)
		return nil, err
	}
	if tlsConfig == nil {
		tlsConfig = defaultTLSConfig.Clone()
	}

	authers = a.autherList(cfg.Handler.Auther, cfg.Handler.Authers...)
	if len(authers) == 0 {
		if auther := a.ParseAutherFromAuth(cfg.Handler.Auth); auther != nil {
			authers = append(authers, auther)
		}
	}

	auther = nil
	if len(authers) > 0 {
		auther = auth.AuthenticatorGroup(authers...)
	}

	var recorders []recorder.RecorderObject
	for _, r := range cfg.Recorders {
		recorders = append(recorders, recorder.RecorderObject{
			Recorder: a.RecorderRegistry().Get(r.Name),
			Record:   r.Record,
		})
	}

	routerOpts := []chain.RouterOption{
		chain.RetriesRouterOption(cfg.Handler.Retries),
		// chain.TimeoutRouterOption(10*time.Second),
		chain.InterfaceRouterOption(ifce),
		chain.SockOptsRouterOption(sockOpts),
		chain.ResolverRouterOption(a.ResolverRegistry().Get(cfg.Resolver)),
		chain.HostMapperRouterOption(a.HostsRegistry().Get(cfg.Hosts)),
		chain.RecordersRouterOption(recorders...),
		chain.LoggerRouterOption(handlerLogger),
	}
	if !ignoreChain {
		routerOpts = append(routerOpts,
			chain.ChainRouterOption(a.chainGroup(cfg.Handler.Chain, cfg.Handler.ChainGroup)),
		)
	}
	router := chain.NewRouter(routerOpts...)

	var h handler.IHandler
	if rf := a.HandlerRegistry().Get(cfg.Handler.Type); rf != nil {
		h = rf(
			handler.RouterOption(router),
			handler.AutherOption(auther),
			handler.AuthOption(a.parseAuth(cfg.Handler.Auth)),
			handler.BypassOption(bypass.BypassGroup(a.bypassList(cfg.Bypass, cfg.Bypasses...)...)),
			handler.TLSConfigOption(tlsConfig),
			handler.RateLimiterOption(a.RateLimiterRegistry().Get(cfg.RLimiter)),
			handler.LoggerOption(handlerLogger),
			handler.ServiceOption(cfg.Name),
		)
	} else {
		return nil, fmt.Errorf("unregistered handler: %s", cfg.Handler.Type)
	}

	if forwarder, ok := h.(handler.IForwarder); ok {
		hop, err := a.parseForwarder(cfg.Forwarder)
		if err != nil {
			return nil, err
		}
		forwarder.Forward(hop)
	}

	if cfg.Handler.Metadata == nil {
		cfg.Handler.Metadata = make(map[string]any)
	}
	handlerLogger.Debugf("metadata: %v", cfg.Handler.Metadata)
	if err := h.Init(metadata.NewMetadata(cfg.Handler.Metadata)); err != nil {
		handlerLogger.Error("init: ", err)
		return nil, err
	}

	s := xservice.NewService(cfg.Name, ln, h,
		xservice.AdmissionOption(admission.AdmissionGroup(admissions...)),
		xservice.PreUpOption(preUp),
		xservice.PreDownOption(preDown),
		xservice.PostUpOption(postUp),
		xservice.PostDownOption(postDown),
		xservice.RecordersOption(recorders...),
		xservice.LoggerOption(serviceLogger),
	)

	serviceLogger.Infof("listening on %s/%s", s.Addr().String(), s.Addr().Network())
	return s, nil
}

func (a *Application) parseForwarder(cfg *config.ForwarderConfig) (chain.IHop, error) {
	if cfg == nil {
		return nil, nil
	}

	hc := config.HopConfig{
		Name:     cfg.Name,
		Selector: cfg.Selector,
	}
	if len(cfg.Nodes) > 0 {
		for _, node := range cfg.Nodes {
			if node != nil {
				hc.Nodes = append(hc.Nodes,
					&config.NodeConfig{
						Name:     node.Name,
						Addr:     node.Addr,
						Host:     node.Host,
						Protocol: node.Protocol,
						Bypass:   node.Bypass,
						Bypasses: node.Bypasses,
						HTTP:     node.HTTP,
						TLS:      node.TLS,
						Auth:     node.Auth,
					},
				)
			}
		}
	} else {
		for _, target := range cfg.Targets {
			if v := strings.TrimSpace(target); v != "" {
				hc.Nodes = append(hc.Nodes,
					&config.NodeConfig{
						Name: target,
						Addr: target,
					},
				)
			}
		}
	}

	if len(hc.Nodes) > 0 {
		return a.ParseHop(&hc)
	}
	return a.HopRegistry().Get(hc.Name), nil
}

func (a *Application) bypassList(name string, names ...string) []bypass.IBypass {
	var bypasses []bypass.IBypass
	if bp := a.BypassRegistry().Get(name); bp != nil {
		bypasses = append(bypasses, bp)
	}
	for _, s := range names {
		if bp := a.BypassRegistry().Get(s); bp != nil {
			bypasses = append(bypasses, bp)
		}
	}
	return bypasses
}

func (a *Application) autherList(name string, names ...string) []auth.IAuthenticator {
	var authers []auth.IAuthenticator
	if auther := a.AutherRegistry().Get(name); auther != nil {
		authers = append(authers, auther)
	}
	for _, s := range names {
		if auther := a.AutherRegistry().Get(s); auther != nil {
			authers = append(authers, auther)
		}
	}
	return authers
}

func (a *Application) admissionList(name string, names ...string) []admission.IAdmission {
	var admissions []admission.IAdmission
	if adm := a.AdmissionRegistry().Get(name); adm != nil {
		admissions = append(admissions, adm)
	}
	for _, s := range names {
		if adm := a.AdmissionRegistry().Get(s); adm != nil {
			admissions = append(admissions, adm)
		}
	}

	return admissions
}

func (a *Application) chainGroup(name string, group *config.ChainGroupConfig) chain.IChainer {
	var chains []chain.IChainer
	var sel selector.ISelector[chain.IChainer]

	if c := a.ChainRegistry().Get(name); c != nil {
		chains = append(chains, c)
	}
	if group != nil {
		for _, s := range group.Chains {
			if c := a.ChainRegistry().Get(s); c != nil {
				chains = append(chains, c)
			}
		}
		sel = a.parseChainSelector(group.Selector)
	}
	if len(chains) == 0 {
		return nil
	}

	if sel == nil {
		sel = a.defaultChainSelector()
	}

	return xchain.NewChainGroup(chains...).
		WithSelector(sel)
}
