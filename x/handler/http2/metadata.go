package http2

import (
	"net/http"
	"strings"
	"time"

	mdata "github.com/jxo-me/netx/core/metadata"
	mdutil "github.com/jxo-me/netx/x/metadata/util"
)

const (
	defaultRealm = "gost"
)

type metadata struct {
	probeResistance        *probeResistance
	header                 http.Header
	hash                   string
	authBasicRealm         string
	observePeriod          time.Duration
	limiterRefreshInterval time.Duration
}

func (h *http2Handler) parseMetadata(md mdata.IMetaData) error {
	if m := mdutil.GetStringMapString(md, "http.header", "header"); len(m) > 0 {
		hd := http.Header{}
		for k, v := range m {
			hd.Add(k, v)
		}
		h.md.header = hd
	}

	if pr := mdutil.GetString(md, "probeResist", "probe_resist"); pr != "" {
		if ss := strings.SplitN(pr, ":", 2); len(ss) == 2 {
			h.md.probeResistance = &probeResistance{
				Type:  ss[0],
				Value: ss[1],
				Knock: mdutil.GetString(md, "knock"),
			}
		}
	}
	h.md.hash = mdutil.GetString(md, "hash")
	h.md.authBasicRealm = mdutil.GetString(md, "authBasicRealm")

	h.md.observePeriod = mdutil.GetDuration(md, "observePeriod", "observer.observePeriod")
	if h.md.observePeriod == 0 {
		h.md.observePeriod = 5 * time.Second
	}
	if h.md.observePeriod < time.Second {
		h.md.observePeriod = time.Second
	}

	h.md.limiterRefreshInterval = mdutil.GetDuration(md, "limiter.refreshInterval")
	if h.md.limiterRefreshInterval == 0 {
		h.md.limiterRefreshInterval = 30 * time.Second
	}
	if h.md.limiterRefreshInterval < time.Second {
		h.md.limiterRefreshInterval = time.Second
	}

	return nil
}

type probeResistance struct {
	Type  string
	Value string
	Knock string
}
