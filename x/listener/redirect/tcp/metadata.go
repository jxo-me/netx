package tcp

import (
	mdata "github.com/jxo-me/netx/core/metadata"
	mdutil "github.com/jxo-me/netx/x/metadata/util"
)

type metadata struct {
	tproxy bool
	mptcp  bool
}

func (l *redirectListener) parseMetadata(md mdata.IMetaData) (err error) {
	l.md.tproxy = mdutil.GetBool(md, "tproxy")
	l.md.mptcp = mdutil.GetBool(md, "mptcp")

	return
}
