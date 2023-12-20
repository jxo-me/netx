package registry

import (
	"github.com/jxo-me/netx/core/logger"
)

type loggerRegistry struct {
	registry[logger.ILogger]
}

func (r *loggerRegistry) Register(name string, v logger.ILogger) error {
	return r.registry.Register(name, v)
}
