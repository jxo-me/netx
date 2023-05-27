package v5

import (
	"time"

	mdata "github.com/jxo-me/netx/core/metadata"
	mdutil "github.com/jxo-me/netx/core/metadata/util"
)

const (
	defaultUDPBufferSize = 4096
)

type metadata struct {
	connectTimeout time.Duration
	noTLS          bool
	relay          string
	udpBufferSize  int
}

func (c *socks5Connector) parseMetadata(md mdata.IMetaData) (err error) {
	const (
		connectTimeout = "timeout"
		noTLS          = "notls"
		relay          = "relay"
		udpBufferSize  = "udpBufferSize"
	)

	c.md.connectTimeout = mdutil.GetDuration(md, connectTimeout)
	c.md.noTLS = mdutil.GetBool(md, noTLS)
	c.md.relay = mdutil.GetString(md, relay)
	c.md.udpBufferSize = mdutil.GetInt(md, udpBufferSize)
	if c.md.udpBufferSize <= 0 {
		c.md.udpBufferSize = defaultUDPBufferSize
	}

	return
}
