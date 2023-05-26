package ss

import (
	"time"

	mdata "github.com/jxo-me/netx/sdk/core/metadata"
	mdutil "github.com/jxo-me/netx/sdk/core/metadata/util"
)

type metadata struct {
	key            string
	connectTimeout time.Duration
	noDelay        bool
}

func (c *ssConnector) parseMetadata(md mdata.IMetaData) (err error) {
	const (
		key            = "key"
		connectTimeout = "timeout"
		noDelay        = "nodelay"
	)

	c.md.key = mdutil.GetString(md, key)
	c.md.connectTimeout = mdutil.GetDuration(md, connectTimeout)
	c.md.noDelay = mdutil.GetBool(md, noDelay)

	return
}
