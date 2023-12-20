package mtls

import (
	"time"

	mdata "github.com/jxo-me/netx/core/metadata"
	mdutil "github.com/jxo-me/netx/core/metadata/util"
	"github.com/jxo-me/netx/x/internal/util/mux"
)

type metadata struct {
	handshakeTimeout time.Duration
	muxCfg           *mux.Config
}

func (d *mtlsDialer) parseMetadata(md mdata.IMetaData) (err error) {
	d.md.handshakeTimeout = mdutil.GetDuration(md, "handshakeTimeout")

	d.md.muxCfg = &mux.Config{
		Version:           mdutil.GetInt(md, "mux.version"),
		KeepAliveInterval: mdutil.GetDuration(md, "mux.keepaliveInterval"),
		KeepAliveDisabled: mdutil.GetBool(md, "mux.keepaliveDisabled"),
		KeepAliveTimeout:  mdutil.GetDuration(md, "mux.keepaliveTimeout"),
		MaxFrameSize:      mdutil.GetInt(md, "mux.maxFrameSize"),
		MaxReceiveBuffer:  mdutil.GetInt(md, "mux.maxReceiveBuffer"),
		MaxStreamBuffer:   mdutil.GetInt(md, "mux.maxStreamBuffer"),
	}
	return
}
