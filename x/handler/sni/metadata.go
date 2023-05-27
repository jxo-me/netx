package sni

import (
	"time"

	mdata "github.com/jxo-me/netx/core/metadata"
	mdutil "github.com/jxo-me/netx/core/metadata/util"
)

type metadata struct {
	readTimeout time.Duration
	hash        string
}

func (h *sniHandler) parseMetadata(md mdata.Metadata) (err error) {
	const (
		readTimeout = "readTimeout"
		hash        = "hash"
	)

	h.md.readTimeout = mdutil.GetDuration(md, readTimeout)
	h.md.hash = mdutil.GetString(md, hash)
	return
}
