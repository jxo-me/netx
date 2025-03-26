package ss

import (
	"time"

	mdata "github.com/jxo-me/netx/core/metadata"
	mdutil "github.com/jxo-me/netx/x/metadata/util"
)

type metadata struct {
	key            string
	connectTimeout time.Duration
}

func (c *ssuConnector) parseMetadata(md mdata.IMetaData) (err error) {
	const (
		key            = "key"
		connectTimeout = "timeout"
	)

	c.md.key = mdutil.GetString(md, key)
	c.md.connectTimeout = mdutil.GetDuration(md, connectTimeout)

	return
}
