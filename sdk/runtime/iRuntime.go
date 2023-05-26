package runtime

import (
	"github.com/jxo-me/netx/sdk/core/admission"
	"github.com/jxo-me/netx/sdk/core/auth"
	"github.com/jxo-me/netx/sdk/core/bypass"
	"github.com/jxo-me/netx/sdk/core/chain"
	"github.com/jxo-me/netx/sdk/core/hosts"
	"github.com/jxo-me/netx/sdk/core/ingress"
	"github.com/jxo-me/netx/sdk/core/limiter/conn"
	"github.com/jxo-me/netx/sdk/core/limiter/rate"
	"github.com/jxo-me/netx/sdk/core/limiter/traffic"
	"github.com/jxo-me/netx/sdk/core/logger"
	"github.com/jxo-me/netx/sdk/core/recorder"
	"github.com/jxo-me/netx/sdk/core/registry"
	"github.com/jxo-me/netx/sdk/core/resolver"
	"github.com/jxo-me/netx/sdk/core/service"

	"github.com/jxo-me/netx/sdk/config"
)

type Runtime interface {
	AdmissionRegistry() registry.IRegistry[admission.IAdmission]
	AutherRegistry() registry.IRegistry[auth.IAuthenticator]
	BypassRegistry() registry.IRegistry[bypass.IBypass]
	ChainRegistry() registry.IRegistry[chain.IChainer]
	ConnectorRegistry() registry.IRegistry[registry.NewConnector]
	ConnLimiterRegistry() registry.IRegistry[conn.IConnLimiter]
	DialerRegistry() registry.IRegistry[registry.NewDialer]
	HandlerRegistry() registry.IRegistry[registry.NewHandler]
	HopRegistry() registry.IRegistry[chain.IHop]
	HostsRegistry() registry.IRegistry[hosts.IHostMapper]
	IngressRegistry() registry.IRegistry[ingress.IIngress]
	ListenerRegistry() registry.IRegistry[registry.NewListener]
	RateLimiterRegistry() registry.IRegistry[rate.IRateLimiter]
	RecorderRegistry() registry.IRegistry[recorder.IRecorder]
	ResolverRegistry() registry.IRegistry[resolver.IResolver]
	ServiceRegistry() registry.IRegistry[service.IService]
	TrafficLimiterRegistry() registry.IRegistry[traffic.ITrafficLimiter]

	// config
	BuildService(cfg *config.Config) (services []service.IService)
	LogFromConfig(cfg *config.LogConfig) logger.ILogger
	BuildAPIService(cfg *config.APIConfig) (service.IService, error)
	BuildMetricsService(cfg *config.MetricsConfig) (service.IService, error)

	ParseChain(cfg *config.ChainConfig) (chain.IChainer, error)
	ParseHop(cfg *config.HopConfig) (chain.IHop, error)
	ParseAuther(cfg *config.AutherConfig) auth.IAuthenticator
	ParseAutherFromAuth(au *config.AuthConfig) auth.IAuthenticator
	ParseAdmission(cfg *config.AdmissionConfig) admission.IAdmission
	ParseBypass(cfg *config.BypassConfig) bypass.IBypass
	ParseResolver(cfg *config.ResolverConfig) (resolver.IResolver, error)
	ParseHosts(cfg *config.HostsConfig) hosts.IHostMapper
	ParseIngress(cfg *config.IngressConfig) ingress.IIngress
	ParseRecorder(cfg *config.RecorderConfig) (r recorder.IRecorder)
	ParseTrafficLimiter(cfg *config.LimiterConfig) (lim traffic.ITrafficLimiter)
	ParseConnLimiter(cfg *config.LimiterConfig) (lim conn.IConnLimiter)
	ParseRateLimiter(cfg *config.LimiterConfig) (lim rate.IRateLimiter)

	ParseService(cfg *config.ServiceConfig) (service.IService, error)

	BuildDefaultTLSConfig(cfg *config.TLSConfig)
}
