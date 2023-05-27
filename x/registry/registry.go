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
	listenerReg  reg.Registry[NewListener]         = new(listenerRegistry)
	handlerReg   reg.Registry[NewHandler]          = new(handlerRegistry)
	dialerReg    reg.Registry[NewDialer]           = new(dialerRegistry)
	connectorReg reg.Registry[NewConnector]        = new(connectorRegistry)
	serviceReg   reg.Registry[service.Service]     = new(serviceRegistry)
	chainReg     reg.Registry[chain.Chainer]       = new(chainRegistry)
	hopReg       reg.Registry[chain.Hop]           = new(hopRegistry)
	autherReg    reg.Registry[auth.Authenticator]  = new(autherRegistry)
	admissionReg reg.Registry[admission.Admission] = new(admissionRegistry)
	bypassReg    reg.Registry[bypass.Bypass]       = new(bypassRegistry)
	resolverReg  reg.Registry[resolver.Resolver]   = new(resolverRegistry)
	hostsReg     reg.Registry[hosts.HostMapper]    = new(hostsRegistry)
	recorderReg  reg.Registry[recorder.Recorder]   = new(recorderRegistry)

	trafficLimiterReg reg.Registry[traffic.ITrafficLimiter] = new(trafficLimiterRegistry)
	connLimiterReg    reg.Registry[conn.IConnLimiter]       = new(connLimiterRegistry)
	rateLimiterReg    reg.Registry[rate.IRateLimiter]       = new(rateLimiterRegistry)

	ingressReg reg.Registry[ingress.Ingress] = new(ingressRegistry)
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

func ListenerRegistry() reg.Registry[NewListener] {
	return listenerReg
}

func HandlerRegistry() reg.Registry[NewHandler] {
	return handlerReg
}

func DialerRegistry() reg.Registry[NewDialer] {
	return dialerReg
}

func ConnectorRegistry() reg.Registry[NewConnector] {
	return connectorReg
}

func ServiceRegistry() reg.Registry[service.Service] {
	return serviceReg
}

func ChainRegistry() reg.Registry[chain.Chainer] {
	return chainReg
}

func HopRegistry() reg.Registry[chain.Hop] {
	return hopReg
}

func AutherRegistry() reg.Registry[auth.Authenticator] {
	return autherReg
}

func AdmissionRegistry() reg.Registry[admission.Admission] {
	return admissionReg
}

func BypassRegistry() reg.Registry[bypass.Bypass] {
	return bypassReg
}

func ResolverRegistry() reg.Registry[resolver.Resolver] {
	return resolverReg
}

func HostsRegistry() reg.Registry[hosts.HostMapper] {
	return hostsReg
}

func RecorderRegistry() reg.Registry[recorder.Recorder] {
	return recorderReg
}

func TrafficLimiterRegistry() reg.Registry[traffic.ITrafficLimiter] {
	return trafficLimiterReg
}

func ConnLimiterRegistry() reg.Registry[conn.IConnLimiter] {
	return connLimiterReg
}

func RateLimiterRegistry() reg.Registry[rate.IRateLimiter] {
	return rateLimiterReg
}

func IngressRegistry() reg.Registry[ingress.Ingress] {
	return ingressReg
}
