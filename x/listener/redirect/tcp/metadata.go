package tcp

import (
	"time"

	mdata "github.com/jxo-me/netx/core/metadata"
	mdutil "github.com/jxo-me/netx/x/metadata/util"
)

type metadata struct {
	tproxy                 bool
	mptcp                  bool
	limiterRefreshInterval time.Duration
}

func (l *redirectListener) parseMetadata(md mdata.IMetaData) (err error) {
	l.md.tproxy = mdutil.GetBool(md, "tproxy")
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
