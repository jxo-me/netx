package udp

import (
	"time"

	md "github.com/jxo-me/netx/sdk/core/metadata"
)

const (
	dialTimeout = "dialTimeout"
)

const (
	defaultDialTimeout = 5 * time.Second
)

type metadata struct {
	dialTimeout time.Duration
}

func (d *udpDialer) parseMetadata(md md.IMetaData) (err error) {
	return
}
