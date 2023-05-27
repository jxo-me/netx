package runtime

import (
	"github.com/jxo-me/netx/core/admission"
	"github.com/jxo-me/netx/core/auth"
	"github.com/jxo-me/netx/core/bypass"
	"github.com/jxo-me/netx/core/chain"
	"github.com/jxo-me/netx/core/hosts"
	"github.com/jxo-me/netx/core/ingress"
	"github.com/jxo-me/netx/core/limiter/conn"
	"github.com/jxo-me/netx/core/limiter/rate"
	"github.com/jxo-me/netx/core/limiter/traffic"
	"github.com/jxo-me/netx/core/logger"
	"github.com/jxo-me/netx/core/recorder"
	reg "github.com/jxo-me/netx/core/registry"
	"github.com/jxo-me/netx/core/resolver"
	"github.com/jxo-me/netx/core/service"
	"github.com/jxo-me/netx/x/registry"

	"github.com/jxo-me/netx/x/config"
)

type IRuntime interface {
	AdmissionRegistry() reg.IRegistry[admission.IAdmission]
	AutherRegistry() reg.IRegistry[auth.IAuthenticator]
	BypassRegistry() reg.IRegistry[bypass.IBypass]
	ChainRegistry() reg.IRegistry[chain.IChainer]
	ConnectorRegistry() reg.IRegistry[registry.NewConnector]
	ConnLimiterRegistry() reg.IRegistry[conn.IConnLimiter]
	DialerRegistry() reg.IRegistry[registry.NewDialer]
	HandlerRegistry() reg.IRegistry[registry.NewHandler]
	HopRegistry() reg.IRegistry[chain.IHop]
	HostsRegistry() reg.IRegistry[hosts.IHostMapper]
	IngressRegistry() reg.IRegistry[ingress.IIngress]
	ListenerRegistry() reg.IRegistry[registry.NewListener]
	RateLimiterRegistry() reg.IRegistry[rate.IRateLimiter]
	RecorderRegistry() reg.IRegistry[recorder.IRecorder]
	ResolverRegistry() reg.IRegistry[resolver.IResolver]
	ServiceRegistry() reg.IRegistry[service.IService]
	TrafficLimiterRegistry() reg.IRegistry[traffic.ITrafficLimiter]

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
