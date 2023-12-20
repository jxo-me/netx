package registry

import (
	"github.com/jxo-me/netx/core/handler"
	"github.com/jxo-me/netx/core/logger"
)

type HandlerRegistry struct {
	registry[handler.NewHandler]
}

func (r *HandlerRegistry) Register(name string, v handler.NewHandler) error {
	if err := r.registry.Register(name, v); err != nil {
		logger.Default().Fatal(err)
	}
	return nil
}
