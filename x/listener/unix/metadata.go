package unix

import (
	md "github.com/jxo-me/netx/core/metadata"
)

type metadata struct{}

func (l *unixListener) parseMetadata(md md.IMetaData) (err error) {
	return
}
