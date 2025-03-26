package quic

import (
	"time"

	mdata "github.com/jxo-me/netx/core/metadata"
	mdutil "github.com/jxo-me/netx/x/metadata/util"
)

const (
	defaultBacklog = 128
)

type metadata struct {
	keepAlivePeriod  time.Duration
	handshakeTimeout time.Duration
	maxIdleTimeout   time.Duration

	backlog int
}

func (l *icmpListener) parseMetadata(md mdata.IMetaData) (err error) {
	l.md.backlog = mdutil.GetInt(md, "backlog")
	if l.md.backlog <= 0 {
		l.md.backlog = defaultBacklog
	}

	if mdutil.GetBool(md, "keepalive") {
		l.md.keepAlivePeriod = mdutil.GetDuration(md, "ttl")
		if l.md.keepAlivePeriod <= 0 {
			l.md.keepAlivePeriod = 10 * time.Second
		}
	}
	l.md.handshakeTimeout = mdutil.GetDuration(md, "handshakeTimeout")
	l.md.maxIdleTimeout = mdutil.GetDuration(md, "maxIdleTimeout")

	return
}
