package v4

import (
	"time"

	mdata "github.com/jxo-me/netx/core/metadata"
	mdutil "github.com/jxo-me/netx/core/metadata/util"
)

type metadata struct {
	readTimeout   time.Duration
	hash          string
	observePeriod time.Duration
}

func (h *socks4Handler) parseMetadata(md mdata.IMetaData) (err error) {
	h.md.readTimeout = mdutil.GetDuration(md, "readTimeout")
	h.md.hash = mdutil.GetString(md, "hash")
	h.md.observePeriod = mdutil.GetDuration(md, "observePeriod")
	return
}
