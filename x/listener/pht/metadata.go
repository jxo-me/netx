package pht

import (
	"strings"
	"time"

	mdata "github.com/jxo-me/netx/core/metadata"
	mdutil "github.com/jxo-me/netx/x/metadata/util"
)

const (
	defaultAuthorizePath = "/authorize"
	defaultPushPath      = "/push"
	defaultPullPath      = "/pull"
	defaultBacklog       = 128
)

type metadata struct {
	authorizePath          string
	pushPath               string
	pullPath               string
	backlog                int
	mptcp                  bool
	limiterRefreshInterval time.Duration
}

func (l *phtListener) parseMetadata(md mdata.IMetaData) (err error) {
	const (
		authorizePath = "authorizePath"
		pushPath      = "pushPath"
		pullPath      = "pullPath"

		backlog = "backlog"
	)

	l.md.authorizePath = mdutil.GetString(md, authorizePath)
	if !strings.HasPrefix(l.md.authorizePath, "/") {
		l.md.authorizePath = defaultAuthorizePath
	}
	l.md.pushPath = mdutil.GetString(md, pushPath)
	if !strings.HasPrefix(l.md.pushPath, "/") {
		l.md.pushPath = defaultPushPath
	}
	l.md.pullPath = mdutil.GetString(md, pullPath)
	if !strings.HasPrefix(l.md.pullPath, "/") {
		l.md.pullPath = defaultPullPath
	}

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
