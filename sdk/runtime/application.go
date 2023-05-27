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
	"github.com/jxo-me/netx/core/recorder"
	reg "github.com/jxo-me/netx/core/registry"
	"github.com/jxo-me/netx/core/resolver"
	"github.com/jxo-me/netx/core/service"
	"github.com/jxo-me/netx/x/registry"
)

type Application struct {
	admissionReg      reg.IRegistry[admission.IAdmission]
	autherReg         reg.IRegistry[auth.IAuthenticator]
	bypassReg         reg.IRegistry[bypass.IBypass]
	chainReg          reg.IRegistry[chain.IChainer]
	connectorReg      reg.IRegistry[registry.NewConnector]
	connLimiterReg    reg.IRegistry[conn.IConnLimiter]
	dialerReg         reg.IRegistry[registry.NewDialer]
	handlerReg        reg.IRegistry[registry.NewHandler]
	hopReg            reg.IRegistry[chain.IHop]
	hostsReg          reg.IRegistry[hosts.IHostMapper]
	ingressReg        reg.IRegistry[ingress.IIngress]
	listenerReg       reg.IRegistry[registry.NewListener]
	rateLimiterReg    reg.IRegistry[rate.IRateLimiter]
	recorderReg       reg.IRegistry[recorder.IRecorder]
	resolverReg       reg.IRegistry[resolver.IResolver]
	serviceReg        reg.IRegistry[service.IService]
	trafficLimiterReg reg.IRegistry[traffic.ITrafficLimiter]
}

func NewConfig() *Application {
	app := Application{
		admissionReg:      new(registry.AdmissionRegistry),
		autherReg:         new(registry.AutherRegistry),
		bypassReg:         new(registry.BypassRegistry),
		chainReg:          new(registry.ChainRegistry),
		connectorReg:      new(registry.ConnectorRegistry),
		connLimiterReg:    new(registry.ConnLimiterRegistry),
		dialerReg:         new(registry.DialerRegistry),
		handlerReg:        new(registry.HandlerRegistry),
		hopReg:            new(registry.HopRegistry),
		hostsReg:          new(registry.HostsRegistry),
		ingressReg:        new(registry.IngressRegistry),
		listenerReg:       new(registry.ListenerRegistry),
		rateLimiterReg:    new(registry.RateLimiterRegistry),
		recorderReg:       new(registry.RecorderRegistry),
		resolverReg:       new(registry.ResolverRegistry),
		serviceReg:        new(registry.ServiceRegistry),
		trafficLimiterReg: new(registry.TrafficLimiterRegistry),
	}

	// Register connectors
	app.InitConnector()
	// Register dialers
	app.InitDialer()
	// Register handlers
	app.InitHandler()
	// Register listeners
	app.InitListener()

	return &app
}

func (a *Application) AdmissionRegistry() reg.IRegistry[admission.IAdmission] {
	return a.admissionReg
}

func (a *Application) AutherRegistry() reg.IRegistry[auth.IAuthenticator] {
	return a.autherReg
}

func (a *Application) BypassRegistry() reg.IRegistry[bypass.IBypass] {
	return a.bypassReg
}

func (a *Application) ChainRegistry() reg.IRegistry[chain.IChainer] {
	return a.chainReg
}

func (a *Application) ConnectorRegistry() reg.IRegistry[registry.NewConnector] {
	return a.connectorReg
}

func (a *Application) ConnLimiterRegistry() reg.IRegistry[conn.IConnLimiter] {
	return a.connLimiterReg
}

func (a *Application) DialerRegistry() reg.IRegistry[registry.NewDialer] {
	return a.dialerReg
}

func (a *Application) HandlerRegistry() reg.IRegistry[registry.NewHandler] {
	return a.handlerReg
}

func (a *Application) HopRegistry() reg.IRegistry[chain.IHop] {
	return a.hopReg
}

func (a *Application) HostsRegistry() reg.IRegistry[hosts.IHostMapper] {
	return a.hostsReg
}

func (a *Application) IngressRegistry() reg.IRegistry[ingress.IIngress] {
	return a.ingressReg
}

func (a *Application) ListenerRegistry() reg.IRegistry[registry.NewListener] {
	return a.listenerReg
}

func (a *Application) RateLimiterRegistry() reg.IRegistry[rate.IRateLimiter] {
	return a.rateLimiterReg
}

func (a *Application) RecorderRegistry() reg.IRegistry[recorder.IRecorder] {
	return a.recorderReg
}

func (a *Application) ResolverRegistry() reg.IRegistry[resolver.IResolver] {
	return a.resolverReg
}

func (a *Application) ServiceRegistry() reg.IRegistry[service.IService] {
	return a.serviceReg
}

func (a *Application) TrafficLimiterRegistry() reg.IRegistry[traffic.ITrafficLimiter] {
	return a.trafficLimiterReg
}
