package tls

import (
	"time"

	md "github.com/jxo-me/netx/core/metadata"
	mdutil "github.com/jxo-me/netx/x/metadata/util"
)

type metadata struct {
	mptcp                  bool
	limiterRefreshInterval time.Duration
}

func (l *obfsListener) parseMetadata(md md.IMetaData) (err error) {
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
