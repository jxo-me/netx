package http

import (
	"net/http"
	"time"

	mdata "github.com/jxo-me/netx/core/metadata"
	mdutil "github.com/jxo-me/netx/x/metadata/util"
)

type metadata struct {
	header                 http.Header
	mptcp                  bool
	limiterRefreshInterval time.Duration
}

func (l *obfsListener) parseMetadata(md mdata.IMetaData) (err error) {
	const (
		header = "header"
	)

	if mm := mdutil.GetStringMapString(md, header); len(mm) > 0 {
		hd := http.Header{}
		for k, v := range mm {
			hd.Add(k, v)
		}
		l.md.header = hd
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
