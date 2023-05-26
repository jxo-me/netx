package registry

import (
	"github.com/jxo-me/netx/sdk/core/connector"
	"github.com/jxo-me/netx/sdk/core/logger"
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
