package registry

import (
	"github.com/jxo-me/netx/core/handler"
	"github.com/jxo-me/netx/core/logger"
)

type NewHandler func(opts ...handler.Option) handler.IHandler

type HandlerRegistry struct {
	registry[NewHandler]
}

func (r *HandlerRegistry) Register(name string, v NewHandler) error {
	if err := r.registry.Register(name, v); err != nil {
		logger.Default().Fatal(err)
	}
	return nil
}
