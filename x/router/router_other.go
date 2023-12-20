//go:build !linux

package router

import (
	"github.com/jxo-me/netx/core/router"
)

func (*localRouter) setSysRoutes(routes ...*router.Route) error {
	return nil
}
