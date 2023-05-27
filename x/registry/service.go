package registry

import (
	"github.com/jxo-me/netx/core/service"
)

type serviceRegistry struct {
	registry[service.Service]
}
