package tcp

import (
	mdata "github.com/jxo-me/netx/sdk/core/metadata"
	mdutil "github.com/jxo-me/netx/sdk/core/metadata/util"
)

type metadata struct {
	tproxy bool
}

func (l *redirectListener) parseMetadata(md mdata.IMetaData) (err error) {
	const (
		tproxy = "tproxy"
	)
	l.md.tproxy = mdutil.GetBool(md, tproxy)
	return
}
