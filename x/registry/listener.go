package registry

import (
	"github.com/jxo-me/netx/core/listener"
	"github.com/jxo-me/netx/core/logger"
)

type NewListener func(opts ...listener.Option) listener.IListener

type listenerRegistry struct {
	registry[NewListener]
}

func (r *listenerRegistry) Register(name string, v NewListener) error {
	if err := r.registry.Register(name, v); err != nil {
		logger.Default().Fatal(err)
	}
	return nil
}
