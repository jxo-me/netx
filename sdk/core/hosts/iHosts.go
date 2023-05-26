package hosts

import (
	"context"
	"net"
)

// IHostMapper is a mapping from hostname to IP.
type IHostMapper interface {
	Lookup(ctx context.Context, network, host string) ([]net.IP, bool)
}
