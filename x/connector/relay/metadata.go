package relay

import (
	"time"

	mdata "github.com/jxo-me/netx/core/metadata"
	mdutil "github.com/jxo-me/netx/core/metadata/util"
	"github.com/jxo-me/netx/x/internal/util/mux"
)

type metadata struct {
	connectTimeout time.Duration
	noDelay        bool
	muxCfg         *mux.Config
}

func (c *relayConnector) parseMetadata(md mdata.IMetaData) (err error) {
	const (
		connectTimeout = "connectTimeout"
		noDelay        = "nodelay"
	)

	c.md.connectTimeout = mdutil.GetDuration(md, connectTimeout)
	c.md.noDelay = mdutil.GetBool(md, noDelay)

	c.md.muxCfg = &mux.Config{
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
