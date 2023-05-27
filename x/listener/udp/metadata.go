package udp

import (
	"time"

	mdata "github.com/jxo-me/netx/core/metadata"
	mdutil "github.com/jxo-me/netx/core/metadata/util"
)

const (
	defaultTTL            = 5 * time.Second
	defaultReadBufferSize = 4096
	defaultReadQueueSize  = 1024
	defaultBacklog        = 128
)

type metadata struct {
	readBufferSize int
	readQueueSize  int
	backlog        int
	keepalive      bool
	ttl            time.Duration
}

func (l *udpListener) parseMetadata(md mdata.IMetaData) (err error) {
	const (
		readBufferSize = "readBufferSize"
		readQueueSize  = "readQueueSize"
		backlog        = "backlog"
		keepAlive      = "keepAlive"
		ttl            = "ttl"
	)

	l.md.ttl = mdutil.GetDuration(md, ttl)
	if l.md.ttl <= 0 {
		l.md.ttl = defaultTTL
	}
	l.md.readBufferSize = mdutil.GetInt(md, readBufferSize)
	if l.md.readBufferSize <= 0 {
		l.md.readBufferSize = defaultReadBufferSize
	}

	l.md.readQueueSize = mdutil.GetInt(md, readQueueSize)
	if l.md.readQueueSize <= 0 {
		l.md.readQueueSize = defaultReadQueueSize
	}

	l.md.backlog = mdutil.GetInt(md, backlog)
	if l.md.backlog <= 0 {
		l.md.backlog = defaultBacklog
	}
	l.md.keepalive = mdutil.GetBool(md, keepAlive)

	return
}
