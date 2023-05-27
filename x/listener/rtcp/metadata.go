package rtcp

import (
	mdata "github.com/jxo-me/netx/core/metadata"
)

type metadata struct{}

func (l *rtcpListener) parseMetadata(md mdata.IMetaData) (err error) {
	return
}
