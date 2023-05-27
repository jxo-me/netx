package v4

import (
	"time"

	mdata "github.com/jxo-me/netx/core/metadata"
	mdutil "github.com/jxo-me/netx/core/metadata/util"
)

type metadata struct {
	connectTimeout time.Duration
	disable4a      bool
}

func (c *socks4Connector) parseMetadata(md mdata.IMetaData) (err error) {
	const (
		connectTimeout = "timeout"
		disable4a      = "disable4a"
	)

	c.md.connectTimeout = mdutil.GetDuration(md, connectTimeout)
	c.md.disable4a = mdutil.GetBool(md, disable4a)

	return
}
