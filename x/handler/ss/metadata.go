package ss

import (
	"time"

	mdata "github.com/jxo-me/netx/core/metadata"
	mdutil "github.com/jxo-me/netx/core/metadata/util"
)

type metadata struct {
	key         string
	readTimeout time.Duration
	hash        string
}

func (h *ssHandler) parseMetadata(md mdata.IMetaData) (err error) {
	const (
		key         = "key"
		readTimeout = "readTimeout"
		hash        = "hash"
	)

	h.md.key = mdutil.GetString(md, key)
	h.md.readTimeout = mdutil.GetDuration(md, readTimeout)
	h.md.hash = mdutil.GetString(md, hash)

	return
}
