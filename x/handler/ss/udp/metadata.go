package ss

import (
	"math"
	"time"

	mdata "github.com/jxo-me/netx/core/metadata"
	mdutil "github.com/jxo-me/netx/x/metadata/util"
)

const (
	MaxMessageSize = math.MaxUint16
)

type metadata struct {
	key         string
	readTimeout time.Duration
}

func (h *ssuHandler) parseMetadata(md mdata.IMetaData) (err error) {
	const (
		key         = "key"
		readTimeout = "readTimeout"
	)

	h.md.key = mdutil.GetString(md, key)
	h.md.readTimeout = mdutil.GetDuration(md, readTimeout)

	return
}
