package direct

import (
	"strings"

	mdata "github.com/jxo-me/netx/core/metadata"
	mdutil "github.com/jxo-me/netx/x/metadata/util"
)

type metadata struct {
	action string
}

func (c *directConnector) parseMetadata(md mdata.IMetaData) (err error) {
	c.md.action = strings.ToLower(mdutil.GetString(md, "action"))
	return
}
