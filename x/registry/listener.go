package registry

import (
	"github.com/jxo-me/netx/core/listener"
	"github.com/jxo-me/netx/core/logger"
)

type ListenerRegistry struct {
	registry[listener.NewListener]
}

func (r *ListenerRegistry) Register(name string, v listener.NewListener) error {
	if err := r.registry.Register(name, v); err != nil {
		logger.Default().Fatal(err)
	}
	return nil
}
