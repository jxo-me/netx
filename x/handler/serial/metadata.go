package serial

import (
	"time"

	mdata "github.com/jxo-me/netx/core/metadata"
	mdutil "github.com/jxo-me/netx/core/metadata/util"
)

const (
	defaultPort     = "COM1"
	defaultBaudRate = 9600
)

type metadata struct {
	timeout time.Duration
}

func (h *serialHandler) parseMetadata(md mdata.IMetaData) (err error) {
	h.md.timeout = mdutil.GetDuration(md, "timeout", "serial.timeout", "handler.serial.timeout")
	return
}
