package registry

import (
	"github.com/jxo-me/netx/sdk/core/dialer"
	"github.com/jxo-me/netx/sdk/core/logger"
)

type NewDialer func(opts ...dialer.Option) dialer.IDialer

type DialerRegistry struct {
	registry[NewDialer]
}

func (r *DialerRegistry) Register(name string, v NewDialer) error {
	if err := r.registry.Register(name, v); err != nil {
		logger.Default().Fatal(err)
	}
	return nil
}
