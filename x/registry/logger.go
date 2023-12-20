package registry

import (
	"github.com/jxo-me/netx/core/logger"
)

type LoggerRegistry struct {
	registry[logger.ILogger]
}

func (r *LoggerRegistry) Register(name string, v logger.ILogger) error {
	return r.registry.Register(name, v)
}
