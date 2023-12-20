package registry

import (
	"github.com/jxo-me/netx/core/connector"
	"github.com/jxo-me/netx/core/logger"
)

type ConnectorRegistry struct {
	registry[connector.NewConnector]
}

func (r *ConnectorRegistry) Register(name string, v connector.NewConnector) error {
	if err := r.registry.Register(name, v); err != nil {
		logger.Default().Fatal(err)
	}
	return nil
}
