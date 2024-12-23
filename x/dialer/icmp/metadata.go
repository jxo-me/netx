package quic

import (
	"time"

	mdata "github.com/jxo-me/netx/core/metadata"
	mdutil "github.com/jxo-me/netx/x/metadata/util"
)

type metadata struct {
	keepAlivePeriod  time.Duration
	maxIdleTimeout   time.Duration
	handshakeTimeout time.Duration
}

func (d *icmpDialer) parseMetadata(md mdata.IMetaData) (err error) {
	if mdutil.GetBool(md, "keepalive") {
		d.md.keepAlivePeriod = mdutil.GetDuration(md, "ttl")
		if d.md.keepAlivePeriod <= 0 {
			d.md.keepAlivePeriod = 10 * time.Second
		}
	}
	d.md.handshakeTimeout = mdutil.GetDuration(md, "handshakeTimeout")
	d.md.maxIdleTimeout = mdutil.GetDuration(md, "maxIdleTimeout")

	return
}
