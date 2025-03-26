package serial

import (
	"time"

	md "github.com/jxo-me/netx/core/metadata"
	mdutil "github.com/jxo-me/netx/x/metadata/util"
)

type metadata struct {
	timeout time.Duration
}

func (l *serialListener) parseMetadata(md md.IMetaData) (err error) {
	l.md.timeout = mdutil.GetDuration(md, "timeout", "serial.timeout", "listener.serial.timeout")
	return
}
