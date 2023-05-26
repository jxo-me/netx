package registry

import (
	"github.com/jxo-me/netx/sdk/core/service"
)

type ServiceRegistry struct {
	registry[service.IService]
}
