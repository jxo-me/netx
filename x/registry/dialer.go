package registry

import (
	"github.com/jxo-me/netx/core/dialer"
	"github.com/jxo-me/netx/core/logger"
)

type NewDialer func(opts ...dialer.Option) dialer.IDialer

type dialerRegistry struct {
	registry[NewDialer]
}

func (r *dialerRegistry) Register(name string, v NewDialer) error {
	if err := r.registry.Register(name, v); err != nil {
		logger.Default().Fatal(err)
	}
	return nil
}
