package dns

import (
	"time"

	mdata "github.com/jxo-me/netx/core/metadata"
	mdutil "github.com/jxo-me/netx/x/metadata/util"
)

const (
	defaultBacklog = 128
)

type metadata struct {
	mode                   string
	readBufferSize         int
	readTimeout            time.Duration
	writeTimeout           time.Duration
	backlog                int
	mptcp                  bool
	limiterRefreshInterval time.Duration
}

func (l *dnsListener) parseMetadata(md mdata.IMetaData) (err error) {
	const (
		backlog        = "backlog"
		mode           = "mode"
		readBufferSize = "readBufferSize"
		readTimeout    = "readTimeout"
		writeTimeout   = "writeTimeout"
	)

	l.md.mode = mdutil.GetString(md, mode)
	l.md.readBufferSize = mdutil.GetInt(md, readBufferSize)
	l.md.readTimeout = mdutil.GetDuration(md, readTimeout)
	l.md.writeTimeout = mdutil.GetDuration(md, writeTimeout)

	l.md.backlog = mdutil.GetInt(md, backlog)
	if l.md.backlog <= 0 {
		l.md.backlog = defaultBacklog
	}
	l.md.mptcp = mdutil.GetBool(md, "mptcp")

	l.md.limiterRefreshInterval = mdutil.GetDuration(md, "limiter.refreshInterval")
	if l.md.limiterRefreshInterval == 0 {
		l.md.limiterRefreshInterval = 30 * time.Second
	}
	if l.md.limiterRefreshInterval < time.Second {
		l.md.limiterRefreshInterval = time.Second
	}

	return
}
