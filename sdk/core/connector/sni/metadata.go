package sni

import (
	"time"

	mdata "github.com/jxo-me/netx/sdk/core/metadata"
	mdutil "github.com/jxo-me/netx/sdk/core/metadata/util"
)

type metadata struct {
	host           string
	connectTimeout time.Duration
}

func (c *sniConnector) parseMetadata(md mdata.IMetaData) (err error) {
	const (
		host           = "host"
		connectTimeout = "timeout"
	)

	c.md.host = mdutil.GetString(md, host)
	c.md.connectTimeout = mdutil.GetDuration(md, connectTimeout)

	return
}
