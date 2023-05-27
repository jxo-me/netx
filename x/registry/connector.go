package registry

import (
	"github.com/jxo-me/netx/core/connector"
	"github.com/jxo-me/netx/core/logger"
)

type NewConnector func(opts ...connector.Option) connector.IConnector

type ConnectorRegistry struct {
	registry[NewConnector]
}

func (r *ConnectorRegistry) Register(name string, v NewConnector) error {
	if err := r.registry.Register(name, v); err != nil {
		logger.Default().Fatal(err)
	}
	return nil
}
