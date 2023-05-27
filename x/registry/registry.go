package registry

import (
	"errors"
	"io"
	"sync"

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
)

var (
	ErrDup = errors.New("registry: duplicate object")
)

var (
	listenerReg  reg.IRegistry[NewListener]          = new(listenerRegistry)
	handlerReg   reg.IRegistry[NewHandler]           = new(handlerRegistry)
	dialerReg    reg.IRegistry[NewDialer]            = new(dialerRegistry)
	connectorReg reg.IRegistry[NewConnector]         = new(connectorRegistry)
	serviceReg   reg.IRegistry[service.IService]     = new(serviceRegistry)
	chainReg     reg.IRegistry[chain.IChainer]       = new(chainRegistry)
	hopReg       reg.IRegistry[chain.IHop]           = new(hopRegistry)
	autherReg    reg.IRegistry[auth.IAuthenticator]  = new(autherRegistry)
	admissionReg reg.IRegistry[admission.IAdmission] = new(admissionRegistry)
	bypassReg    reg.IRegistry[bypass.IBypass]       = new(bypassRegistry)
	resolverReg  reg.IRegistry[resolver.IResolver]   = new(resolverRegistry)
	hostsReg     reg.IRegistry[hosts.IHostMapper]    = new(hostsRegistry)
	recorderReg  reg.IRegistry[recorder.IRecorder]   = new(recorderRegistry)

	trafficLimiterReg reg.IRegistry[traffic.ITrafficLimiter] = new(trafficLimiterRegistry)
	connLimiterReg    reg.IRegistry[conn.IConnLimiter]       = new(connLimiterRegistry)
	rateLimiterReg    reg.IRegistry[rate.IRateLimiter]       = new(rateLimiterRegistry)

	ingressReg reg.IRegistry[ingress.IIngress] = new(ingressRegistry)
)

type registry[T any] struct {
	m sync.Map
}

func (r *registry[T]) Register(name string, v T) error {
	if name == "" {
		return nil
	}
	if _, loaded := r.m.LoadOrStore(name, v); loaded {
		return ErrDup
	}

	return nil
}

func (r *registry[T]) Unregister(name string) {
	if v, ok := r.m.Load(name); ok {
		if closer, ok := v.(io.Closer); ok {
			closer.Close()
		}
		r.m.Delete(name)
	}
}

func (r *registry[T]) IsRegistered(name string) bool {
	_, ok := r.m.Load(name)
	return ok
}

func (r *registry[T]) Get(name string) (t T) {
	if name == "" {
		return
	}
	v, _ := r.m.Load(name)
	t, _ = v.(T)
	return
}

func (r *registry[T]) GetAll() (m map[string]T) {
	m = make(map[string]T)
	r.m.Range(func(key, value any) bool {
		k, _ := key.(string)
		v, _ := value.(T)
		m[k] = v
		return true
	})
	return
}

func ListenerRegistry() reg.IRegistry[NewListener] {
	return listenerReg
}

func HandlerRegistry() reg.IRegistry[NewHandler] {
	return handlerReg
}

func DialerRegistry() reg.IRegistry[NewDialer] {
	return dialerReg
}

func ConnectorRegistry() reg.IRegistry[NewConnector] {
	return connectorReg
}

func ServiceRegistry() reg.IRegistry[service.IService] {
	return serviceReg
}

func ChainRegistry() reg.IRegistry[chain.IChainer] {
	return chainReg
}

func HopRegistry() reg.IRegistry[chain.IHop] {
	return hopReg
}

func AutherRegistry() reg.IRegistry[auth.IAuthenticator] {
	return autherReg
}

func AdmissionRegistry() reg.IRegistry[admission.IAdmission] {
	return admissionReg
}

func BypassRegistry() reg.IRegistry[bypass.IBypass] {
	return bypassReg
}

func ResolverRegistry() reg.IRegistry[resolver.IResolver] {
	return resolverReg
}

func HostsRegistry() reg.IRegistry[hosts.IHostMapper] {
	return hostsReg
}

func RecorderRegistry() reg.IRegistry[recorder.IRecorder] {
	return recorderReg
}

func TrafficLimiterRegistry() reg.IRegistry[traffic.ITrafficLimiter] {
	return trafficLimiterReg
}

func ConnLimiterRegistry() reg.IRegistry[conn.IConnLimiter] {
	return connLimiterReg
}

func RateLimiterRegistry() reg.IRegistry[rate.IRateLimiter] {
	return rateLimiterReg
}

func IngressRegistry() reg.IRegistry[ingress.IIngress] {
	return ingressReg
}
