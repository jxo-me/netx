package mtls

import (
	"time"

	mdata "github.com/jxo-me/netx/core/metadata"
	mdutil "github.com/jxo-me/netx/core/metadata/util"
)

const (
	defaultBacklog = 128
)

type metadata struct {
	muxKeepAliveDisabled bool
	muxKeepAliveInterval time.Duration
	muxKeepAliveTimeout  time.Duration
	muxMaxFrameSize      int
	muxMaxReceiveBuffer  int
	muxMaxStreamBuffer   int

	backlog int
}

func (l *mtlsListener) parseMetadata(md mdata.IMetaData) (err error) {
	const (
		backlog = "backlog"

		muxKeepAliveDisabled = "muxKeepAliveDisabled"
		muxKeepAliveInterval = "muxKeepAliveInterval"
		muxKeepAliveTimeout  = "muxKeepAliveTimeout"
		muxMaxFrameSize      = "muxMaxFrameSize"
		muxMaxReceiveBuffer  = "muxMaxReceiveBuffer"
		muxMaxStreamBuffer   = "muxMaxStreamBuffer"
	)

	l.md.backlog = mdutil.GetInt(md, backlog)
	if l.md.backlog <= 0 {
		l.md.backlog = defaultBacklog
	}

	l.md.muxKeepAliveDisabled = mdutil.GetBool(md, muxKeepAliveDisabled)
	l.md.muxKeepAliveInterval = mdutil.GetDuration(md, muxKeepAliveInterval)
	l.md.muxKeepAliveTimeout = mdutil.GetDuration(md, muxKeepAliveTimeout)
	l.md.muxMaxFrameSize = mdutil.GetInt(md, muxMaxFrameSize)
	l.md.muxMaxReceiveBuffer = mdutil.GetInt(md, muxMaxReceiveBuffer)
	l.md.muxMaxStreamBuffer = mdutil.GetInt(md, muxMaxStreamBuffer)

	return
}
