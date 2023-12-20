package registry

import (
	"github.com/jxo-me/netx/core/logger"
)

type ListenerRegistry struct {
	registry[NewListener]
}

func (r *ListenerRegistry) Register(name string, v NewListener) error {
	if err := r.registry.Register(name, v); err != nil {
		logger.Default().Fatal(err)
	}
	return nil
}
