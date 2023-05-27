package registry

import (
	"github.com/jxo-me/netx/core/service"
)

type ServiceRegistry struct {
	registry[service.IService]
}
