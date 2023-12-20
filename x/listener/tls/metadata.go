package tls

import (
	mdata "github.com/jxo-me/netx/core/metadata"
	mdutil "github.com/jxo-me/netx/core/metadata/util"
)

type metadata struct {
	mptcp bool
}

func (l *tlsListener) parseMetadata(md mdata.IMetaData) (err error) {
	l.md.mptcp = mdutil.GetBool(md, "mptcp")
	return
}
