package app

import (
	"github.com/jxo-me/netx/core/admission"
	"github.com/jxo-me/netx/core/api"
	"github.com/jxo-me/netx/core/app"
	"github.com/jxo-me/netx/core/auth"
	"github.com/jxo-me/netx/core/bypass"
	"github.com/jxo-me/netx/core/chain"
	"github.com/jxo-me/netx/core/connector"
	"github.com/jxo-me/netx/core/dialer"
	"github.com/jxo-me/netx/core/handler"
	"github.com/jxo-me/netx/core/hop"
	"github.com/jxo-me/netx/core/hosts"
	"github.com/jxo-me/netx/core/ingress"
	"github.com/jxo-me/netx/core/limiter/conn"
	"github.com/jxo-me/netx/core/limiter/rate"
	"github.com/jxo-me/netx/core/limiter/traffic"
	"github.com/jxo-me/netx/core/listener"
	"github.com/jxo-me/netx/core/logger"
	"github.com/jxo-me/netx/core/observer"
	"github.com/jxo-me/netx/core/recorder"
	reg "github.com/jxo-me/netx/core/registry"
	"github.com/jxo-me/netx/core/resolver"
	"github.com/jxo-me/netx/core/router"
	"github.com/jxo-me/netx/core/sd"
	"github.com/jxo-me/netx/core/service"
	"github.com/jxo-me/netx/x/registry"
)

var (
	Runtime app.IRuntime = NewConfig()
	ApiSrv  api.IApi
)

type Application struct {
	admissionReg      reg.IRegistry[admission.IAdmission]
	autherReg         reg.IRegistry[auth.IAuthenticator]
	bypassReg         reg.IRegistry[bypass.IBypass]
	chainReg          reg.IRegistry[chain.IChainer]
	connectorReg      reg.IRegistry[connector.NewConnector]
	connLimiterReg    reg.IRegistry[conn.IConnLimiter]
	dialerReg         reg.IRegistry[dialer.NewDialer]
	handlerReg        reg.IRegistry[handler.NewHandler]
	hopReg            reg.IRegistry[hop.IHop]
	hostsReg          reg.IRegistry[hosts.IHostMapper]
	ingressReg        reg.IRegistry[ingress.IIngress]
	listenerReg       reg.IRegistry[listener.NewListener]
	rateLimiterReg    reg.IRegistry[rate.IRateLimiter]
	recorderReg       reg.IRegistry[recorder.IRecorder]
	resolverReg       reg.IRegistry[resolver.IResolver]
	routerReg         reg.IRegistry[router.IRouter]
	sdReg             reg.IRegistry[sd.ISD]
	observerReg       reg.IRegistry[observer.IObserver]
	serviceReg        reg.IRegistry[service.IService]
	loggerReg         reg.IRegistry[logger.ILogger]
	trafficLimiterReg reg.IRegistry[traffic.ITrafficLimiter]
}

func NewConfig() *Application {
	a := Application{
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
		routerReg:         new(registry.RouterRegistry),
		sdReg:             new(registry.SdRegistry),
		observerReg:       new(registry.ObserverRegistry),
		serviceReg:        new(registry.ServiceRegistry),
		loggerReg:         new(registry.LoggerRegistry),
		trafficLimiterReg: new(registry.TrafficLimiterRegistry),
	}

	// Register connectors
	// Register dialers
	// Register handlers
	// Register listeners

	return &a
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

func (a *Application) ConnectorRegistry() reg.IRegistry[connector.NewConnector] {
	return a.connectorReg
}

func (a *Application) ConnLimiterRegistry() reg.IRegistry[conn.IConnLimiter] {
	return a.connLimiterReg
}

func (a *Application) DialerRegistry() reg.IRegistry[dialer.NewDialer] {
	return a.dialerReg
}

func (a *Application) HandlerRegistry() reg.IRegistry[handler.NewHandler] {
	return a.handlerReg
}

func (a *Application) HopRegistry() reg.IRegistry[hop.IHop] {
	return a.hopReg
}

func (a *Application) HostsRegistry() reg.IRegistry[hosts.IHostMapper] {
	return a.hostsReg
}

func (a *Application) IngressRegistry() reg.IRegistry[ingress.IIngress] {
	return a.ingressReg
}

func (a *Application) ListenerRegistry() reg.IRegistry[listener.NewListener] {
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

func (a *Application) RouterRegistry() reg.IRegistry[router.IRouter] {
	return a.routerReg
}

func (a *Application) SDRegistry() reg.IRegistry[sd.ISD] {
	return a.sdReg
}

func (a *Application) ObserverRegistry() reg.IRegistry[observer.IObserver] {
	return a.observerReg
}

func (a *Application) ServiceRegistry() reg.IRegistry[service.IService] {
	return a.serviceReg
}

func (a *Application) LoggerRegistry() reg.IRegistry[logger.ILogger] {
	return a.loggerReg
}

func (a *Application) TrafficLimiterRegistry() reg.IRegistry[traffic.ITrafficLimiter] {
	return a.trafficLimiterReg
}
