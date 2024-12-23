package service

import (
	"fmt"
	"github.com/jxo-me/netx/x/app"
	"runtime"
	"strings"
	"time"

	"github.com/jxo-me/netx/core/auth"
	"github.com/jxo-me/netx/core/chain"
	"github.com/jxo-me/netx/core/handler"
	"github.com/jxo-me/netx/core/hop"
	"github.com/jxo-me/netx/core/listener"
	"github.com/jxo-me/netx/core/logger"
	"github.com/jxo-me/netx/core/observer/stats"
	"github.com/jxo-me/netx/core/recorder"
	"github.com/jxo-me/netx/core/selector"
	"github.com/jxo-me/netx/core/service"
	xadmission "github.com/jxo-me/netx/x/admission"
	xauth "github.com/jxo-me/netx/x/auth"
	xbypass "github.com/jxo-me/netx/x/bypass"
	xchain "github.com/jxo-me/netx/x/chain"
	"github.com/jxo-me/netx/x/config"
	"github.com/jxo-me/netx/x/config/parsing"
	admission_parser "github.com/jxo-me/netx/x/config/parsing/admission"
	auth_parser "github.com/jxo-me/netx/x/config/parsing/auth"
	bypass_parser "github.com/jxo-me/netx/x/config/parsing/bypass"
	hop_parser "github.com/jxo-me/netx/x/config/parsing/hop"
	logger_parser "github.com/jxo-me/netx/x/config/parsing/logger"
	selector_parser "github.com/jxo-me/netx/x/config/parsing/selector"
	tls_util "github.com/jxo-me/netx/x/internal/util/tls"
	"github.com/jxo-me/netx/x/metadata"
	mdutil "github.com/jxo-me/netx/x/metadata/util"
	xservice "github.com/jxo-me/netx/x/service"
	"github.com/vishvananda/netns"
)

func ParseService(cfg *config.ServiceConfig) (service.IService, error) {
	if cfg.Listener == nil {
		cfg.Listener = &config.ListenerConfig{}
	}
	if strings.TrimSpace(cfg.Listener.Type) == "" {
		cfg.Listener.Type = "tcp"
	}

	if cfg.Handler == nil {
		cfg.Handler = &config.HandlerConfig{}
	}
	if strings.TrimSpace(cfg.Handler.Type) == "" {
		cfg.Handler.Type = "auto"
	}

	log := logger.Default()
	if loggers := logger_parser.List(cfg.Logger, cfg.Loggers...); len(loggers) > 0 {
		log = logger.LoggerGroup(loggers...)
	}

	serviceLogger := log.WithFields(map[string]any{
		"kind":     "service",
		"service":  cfg.Name,
		"listener": cfg.Listener.Type,
		"handler":  cfg.Handler.Type,
	})

	tlsCfg := cfg.Listener.TLS
	if tlsCfg == nil {
		tlsCfg = &config.TLSConfig{}
	}
	tlsConfig, err := tls_util.LoadServerConfig(tlsCfg)
	if err != nil {
		serviceLogger.Error(err)
		return nil, err
	}
	if tlsConfig == nil {
		tlsConfig = parsing.DefaultTLSConfig().Clone()
		tls_util.SetTLSOptions(tlsConfig, tlsCfg.Options)
	}

	authers := auth_parser.List(cfg.Listener.Auther, cfg.Listener.Authers...)
	if len(authers) == 0 {
		if auther := auth_parser.ParseAutherFromAuth(cfg.Listener.Auth); auther != nil {
			authers = append(authers, auther)
		}
	}
	var auther auth.IAuthenticator
	if len(authers) > 0 {
		auther = xauth.AuthenticatorGroup(authers...)
	}

	admissions := admission_parser.List(cfg.Admission, cfg.Admissions...)

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
	var pStats *stats.Stats
	var observePeriod time.Duration
	var netnsIn, netnsOut string
	var dialTimeout time.Duration
	if cfg.Metadata != nil {
		md := metadata.NewMetadata(cfg.Metadata)
		ppv = mdutil.GetInt(md, parsing.MDKeyProxyProtocol)
		if v := mdutil.GetString(md, parsing.MDKeyInterface); v != "" {
			ifce = v
		}
		if v := mdutil.GetInt(md, parsing.MDKeySoMark); v > 0 {
			sockOpts = &chain.SockOpts{
				Mark: v,
			}
		}
		preUp = mdutil.GetStrings(md, parsing.MDKeyPreUp)
		preDown = mdutil.GetStrings(md, parsing.MDKeyPreDown)
		postUp = mdutil.GetStrings(md, parsing.MDKeyPostUp)
		postDown = mdutil.GetStrings(md, parsing.MDKeyPostDown)
		ignoreChain = mdutil.GetBool(md, parsing.MDKeyIgnoreChain)

		if mdutil.GetBool(md, parsing.MDKeyEnableStats) {
			pStats = &stats.Stats{}
		}
		observePeriod = mdutil.GetDuration(md, "observePeriod")
		netnsIn = mdutil.GetString(md, "netns")
		netnsOut = mdutil.GetString(md, "netns.out")
		dialTimeout = mdutil.GetDuration(md, "dialTimeout")
	}

	listenerLogger := serviceLogger.WithFields(map[string]any{
		"kind": "listener",
	})

	routerOpts := []chain.RouterOption{
		chain.TimeoutRouterOption(dialTimeout),
		chain.InterfaceRouterOption(ifce),
		chain.NetnsRouterOption(netnsOut),
		chain.SockOptsRouterOption(sockOpts),
		chain.ResolverRouterOption(app.Runtime.ResolverRegistry().Get(cfg.Resolver)),
		chain.HostMapperRouterOption(app.Runtime.HostsRegistry().Get(cfg.Hosts)),
		chain.LoggerRouterOption(listenerLogger),
	}
	if !ignoreChain {
		routerOpts = append(routerOpts,
			chain.ChainRouterOption(chainGroup(cfg.Listener.Chain, cfg.Listener.ChainGroup)),
		)
	}

	listenOpts := []listener.Option{
		listener.AddrOption(cfg.Addr),
		listener.RouterOption(xchain.NewRouter(routerOpts...)),
		listener.AutherOption(auther),
		listener.AuthOption(auth_parser.Info(cfg.Listener.Auth)),
		listener.TLSConfigOption(tlsConfig),
		listener.AdmissionOption(xadmission.AdmissionGroup(admissions...)),
		listener.TrafficLimiterOption(app.Runtime.TrafficLimiterRegistry().Get(cfg.Limiter)),
		listener.ConnLimiterOption(app.Runtime.ConnLimiterRegistry().Get(cfg.CLimiter)),
		listener.ServiceOption(cfg.Name),
		listener.ProxyProtocolOption(ppv),
		listener.StatsOption(pStats),
		listener.NetnsOption(netnsIn),
		listener.LoggerOption(listenerLogger),
	}

	if netnsIn != "" {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		originNs, err := netns.Get()
		if err != nil {
			return nil, fmt.Errorf("netns.Get(): %v", err)
		}
		defer netns.Set(originNs)

		var ns netns.NsHandle

		if strings.HasPrefix(netnsIn, "/") {
			ns, err = netns.GetFromPath(netnsIn)
		} else {
			ns, err = netns.GetFromName(netnsIn)
		}
		if err != nil {
			return nil, fmt.Errorf("netns.Get(%s): %v", netnsIn, err)
		}
		defer ns.Close()

		if err := netns.Set(ns); err != nil {
			return nil, fmt.Errorf("netns.Set(%s): %v", netnsIn, err)
		}
	}

	var ln listener.IListener
	if rf := app.Runtime.ListenerRegistry().Get(cfg.Listener.Type); rf != nil {
		ln = rf(listenOpts...)
	} else {
		return nil, fmt.Errorf("unknown listener: %s", cfg.Listener.Type)
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
	tlsConfig, err = tls_util.LoadServerConfig(tlsCfg)
	if err != nil {
		handlerLogger.Error(err)
		return nil, err
	}
	if tlsConfig == nil {
		tlsConfig = parsing.DefaultTLSConfig().Clone()
		tls_util.SetTLSOptions(tlsConfig, tlsCfg.Options)
	}

	authers = auth_parser.List(cfg.Handler.Auther, cfg.Handler.Authers...)
	if len(authers) == 0 {
		if auther := auth_parser.ParseAutherFromAuth(cfg.Handler.Auth); auther != nil {
			authers = append(authers, auther)
		}
	}

	auther = nil
	if len(authers) > 0 {
		auther = xauth.AuthenticatorGroup(authers...)
	}

	var recorders []recorder.RecorderObject
	for _, r := range cfg.Recorders {
		md := metadata.NewMetadata(r.Metadata)
		recorders = append(recorders, recorder.RecorderObject{
			Recorder: app.Runtime.RecorderRegistry().Get(r.Name),
			Record:   r.Record,
			Options: &recorder.Options{
				Direction:       mdutil.GetBool(md, parsing.MDKeyRecorderDirection),
				TimestampFormat: mdutil.GetString(md, parsing.MDKeyRecorderTimestampFormat),
				Hexdump:         mdutil.GetBool(md, parsing.MDKeyRecorderHexdump),
				HTTPBody:        mdutil.GetBool(md, parsing.MDKeyRecorderHTTPBody),
				MaxBodySize:     mdutil.GetInt(md, parsing.MDKeyRecorderHTTPMaxBodySize),
			},
		})
	}

	routerOpts = []chain.RouterOption{
		chain.RetriesRouterOption(cfg.Handler.Retries),
		chain.TimeoutRouterOption(dialTimeout),
		chain.InterfaceRouterOption(ifce),
		chain.NetnsRouterOption(netnsOut),
		chain.SockOptsRouterOption(sockOpts),
		chain.ResolverRouterOption(app.Runtime.ResolverRegistry().Get(cfg.Resolver)),
		chain.HostMapperRouterOption(app.Runtime.HostsRegistry().Get(cfg.Hosts)),
		chain.RecordersRouterOption(recorders...),
		chain.LoggerRouterOption(handlerLogger),
	}
	if !ignoreChain {
		routerOpts = append(routerOpts,
			chain.ChainRouterOption(chainGroup(cfg.Handler.Chain, cfg.Handler.ChainGroup)),
		)
	}

	var h handler.IHandler
	if rf := app.Runtime.HandlerRegistry().Get(cfg.Handler.Type); rf != nil {
		h = rf(
			handler.RouterOption(xchain.NewRouter(routerOpts...)),
			handler.AutherOption(auther),
			handler.AuthOption(auth_parser.Info(cfg.Handler.Auth)),
			handler.BypassOption(xbypass.BypassGroup(bypass_parser.List(cfg.Bypass, cfg.Bypasses...)...)),
			handler.TLSConfigOption(tlsConfig),
			handler.RateLimiterOption(app.Runtime.RateLimiterRegistry().Get(cfg.RLimiter)),
			handler.TrafficLimiterOption(app.Runtime.TrafficLimiterRegistry().Get(cfg.Handler.Limiter)),
			handler.ObserverOption(app.Runtime.ObserverRegistry().Get(cfg.Handler.Observer)),
			handler.RecordersOption(recorders...),
			handler.LoggerOption(handlerLogger),
			handler.ServiceOption(cfg.Name),
			handler.NetnsOption(netnsIn),
		)
	} else {
		return nil, fmt.Errorf("unknown handler: %s", cfg.Handler.Type)
	}

	if forwarder, ok := h.(handler.IForwarder); ok {
		hop, err := parseForwarder(cfg.Forwarder, log)
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
		xservice.AdmissionOption(xadmission.AdmissionGroup(admissions...)),
		xservice.PreUpOption(preUp),
		xservice.PreDownOption(preDown),
		xservice.PostUpOption(postUp),
		xservice.PostDownOption(postDown),
		xservice.RecordersOption(recorders...),
		xservice.StatsOption(pStats),
		xservice.ObserverOption(app.Runtime.ObserverRegistry().Get(cfg.Observer)),
		xservice.ObservePeriodOption(observePeriod),
		xservice.LoggerOption(serviceLogger),
	)

	serviceLogger.Infof("listening on %s/%s", s.Addr().String(), s.Addr().Network())
	return s, nil
}

func parseForwarder(cfg *config.ForwarderConfig, log logger.ILogger) (hop.IHop, error) {
	if cfg == nil {
		return nil, nil
	}

	hopName := cfg.Hop
	if hopName == "" {
		hopName = cfg.Name
	}
	if hopName != "" {
		return app.Runtime.HopRegistry().Get(hopName), nil
	}

	hc := config.HopConfig{
		Name:     cfg.Name,
		Selector: cfg.Selector,
	}
	for _, node := range cfg.Nodes {
		if node == nil {
			continue
		}

		filter := node.Filter
		if filter == nil {
			if node.Protocol != "" || node.Host != "" || node.Path != "" {
				filter = &config.NodeFilterConfig{
					Protocol: node.Protocol,
					Host:     node.Host,
					Path:     node.Path,
				}
			}
		}

		httpCfg := node.HTTP
		if node.Auth != nil {
			if httpCfg == nil {
				httpCfg = &config.HTTPNodeConfig{}
			}
			if httpCfg.Auth == nil {
				httpCfg.Auth = node.Auth
			}
		}
		hc.Nodes = append(hc.Nodes, &config.NodeConfig{
			Name:     node.Name,
			Addr:     node.Addr,
			Network:  node.Network,
			Bypass:   node.Bypass,
			Bypasses: node.Bypasses,
			Filter:   filter,
			Matcher:  node.Matcher,
			HTTP:     httpCfg,
			TLS:      node.TLS,
			Metadata: node.Metadata,
		})
	}
	return hop_parser.ParseHop(&hc, log)
}

func chainGroup(name string, group *config.ChainGroupConfig) chain.IChainer {
	var chains []chain.IChainer
	var sel selector.ISelector[chain.IChainer]

	if c := app.Runtime.ChainRegistry().Get(name); c != nil {
		chains = append(chains, c)
	}
	if group != nil {
		for _, s := range group.Chains {
			if c := app.Runtime.ChainRegistry().Get(s); c != nil {
				chains = append(chains, c)
			}
		}
		sel = selector_parser.ParseChainSelector(group.Selector)
	}
	if len(chains) == 0 {
		return nil
	}

	if sel == nil {
		sel = selector_parser.DefaultChainSelector()
	}

	return xchain.NewChainGroup(chains...).
		WithSelector(sel)
}
