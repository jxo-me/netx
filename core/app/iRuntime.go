package app

import (
	"github.com/jxo-me/netx/core/admission"
	"github.com/jxo-me/netx/core/auth"
	"github.com/jxo-me/netx/core/bypass"
	"github.com/jxo-me/netx/core/chain"
	"github.com/jxo-me/netx/core/connector"
	"github.com/jxo-me/netx/core/dialer"
	"github.com/jxo-me/netx/core/handler"
	"github.com/jxo-me/netx/core/hosts"
	"github.com/jxo-me/netx/core/ingress"
	"github.com/jxo-me/netx/core/limiter/conn"
	"github.com/jxo-me/netx/core/limiter/rate"
	"github.com/jxo-me/netx/core/limiter/traffic"
	"github.com/jxo-me/netx/core/listener"
	"github.com/jxo-me/netx/core/recorder"
	reg "github.com/jxo-me/netx/core/registry"
	"github.com/jxo-me/netx/core/resolver"
	"github.com/jxo-me/netx/core/service"
)

type IRuntime interface {
	AdmissionRegistry() reg.IRegistry[admission.IAdmission]
	AutherRegistry() reg.IRegistry[auth.IAuthenticator]
	BypassRegistry() reg.IRegistry[bypass.IBypass]
	ChainRegistry() reg.IRegistry[chain.IChainer]
	ConnectorRegistry() reg.IRegistry[connector.NewConnector]
	ConnLimiterRegistry() reg.IRegistry[conn.IConnLimiter]
	DialerRegistry() reg.IRegistry[dialer.NewDialer]
	HandlerRegistry() reg.IRegistry[handler.NewHandler]
	HopRegistry() reg.IRegistry[chain.IHop]
	HostsRegistry() reg.IRegistry[hosts.IHostMapper]
	IngressRegistry() reg.IRegistry[ingress.IIngress]
	ListenerRegistry() reg.IRegistry[listener.NewListener]
	RateLimiterRegistry() reg.IRegistry[rate.IRateLimiter]
	RecorderRegistry() reg.IRegistry[recorder.IRecorder]
	ResolverRegistry() reg.IRegistry[resolver.IResolver]
	ServiceRegistry() reg.IRegistry[service.IService]
	TrafficLimiterRegistry() reg.IRegistry[traffic.ITrafficLimiter]
}
