package tap

import (
	"math"

	mdata "github.com/jxo-me/netx/core/metadata"
	mdutil "github.com/jxo-me/netx/x/metadata/util"
)

const (
	MaxMessageSize = math.MaxUint16
)

type metadata struct {
	key string
}

func (h *tapHandler) parseMetadata(md mdata.IMetaData) (err error) {
	const (
		key = "key"
	)

	h.md.key = mdutil.GetString(md, key)
	return
}
