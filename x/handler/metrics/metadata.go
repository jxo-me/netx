package metrics

import (
	mdata "github.com/jxo-me/netx/core/metadata"
	mdutil "github.com/jxo-me/netx/core/metadata/util"
)

const (
	DefaultPath = "/metrics"
)

type metadata struct {
	path string
}

func (h *metricsHandler) parseMetadata(md mdata.IMetaData) (err error) {
	h.md.path = mdutil.GetString(md, "metrics.path", "path")
	if h.md.path == "" {
		h.md.path = DefaultPath
	}
	return
}
