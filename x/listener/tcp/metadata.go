package tcp

import (
	md "github.com/jxo-me/netx/core/metadata"
	mdutil "github.com/jxo-me/netx/core/metadata/util"
)

type metadata struct {
	mptcp bool
}

func (l *tcpListener) parseMetadata(md md.IMetaData) (err error) {
	l.md.mptcp = mdutil.GetBool(md, "mptcp")
	return
}
