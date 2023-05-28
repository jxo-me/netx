package registry

import (
	"github.com/jxo-me/netx/core/dialer"
	"github.com/jxo-me/netx/core/logger"
)

type DialerRegistry struct {
	registry[dialer.NewDialer]
}

func (r *DialerRegistry) Register(name string, v dialer.NewDialer) error {
	if err := r.registry.Register(name, v); err != nil {
		logger.Default().Fatal(err)
	}
	return nil
}
