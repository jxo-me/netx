package mws

import (
	"net/http"
	"time"

	mdata "github.com/jxo-me/netx/core/metadata"
	"github.com/jxo-me/netx/x/internal/util/mux"
	mdutil "github.com/jxo-me/netx/x/metadata/util"
)

const (
	defaultPath    = "/ws"
	defaultBacklog = 128
)

type metadata struct {
	path    string
	backlog int
	header  http.Header

	handshakeTimeout  time.Duration
	readHeaderTimeout time.Duration
	readBufferSize    int
	writeBufferSize   int
	enableCompression bool

	muxCfg *mux.Config

	mptcp                  bool
	limiterRefreshInterval time.Duration
}

func (l *mwsListener) parseMetadata(md mdata.IMetaData) (err error) {
	l.md.path = mdutil.GetString(md, "ws.path", "path")
	if l.md.path == "" {
		l.md.path = defaultPath
	}

	l.md.backlog = mdutil.GetInt(md, "ws.backlog", "backlog")
	if l.md.backlog <= 0 {
		l.md.backlog = defaultBacklog
	}

	l.md.handshakeTimeout = mdutil.GetDuration(md, "ws.handshakeTimeout", "handshakeTimeout")
	l.md.readHeaderTimeout = mdutil.GetDuration(md, "ws.readHeaderTimeout", "readHeaderTimeout")
	l.md.readBufferSize = mdutil.GetInt(md, "ws.readBufferSize", "readBufferSize")
	l.md.writeBufferSize = mdutil.GetInt(md, "ws.writeBufferSize", "writeBufferSize")
	l.md.enableCompression = mdutil.GetBool(md, "ws.enableCompression", "enableCompression")

	l.md.muxCfg = &mux.Config{
		Version:           mdutil.GetInt(md, "mux.version"),
		KeepAliveInterval: mdutil.GetDuration(md, "mux.keepaliveInterval"),
		KeepAliveDisabled: mdutil.GetBool(md, "mux.keepaliveDisabled"),
		KeepAliveTimeout:  mdutil.GetDuration(md, "mux.keepaliveTimeout"),
		MaxFrameSize:      mdutil.GetInt(md, "mux.maxFrameSize"),
		MaxReceiveBuffer:  mdutil.GetInt(md, "mux.maxReceiveBuffer"),
		MaxStreamBuffer:   mdutil.GetInt(md, "mux.maxStreamBuffer"),
	}

	if mm := mdutil.GetStringMapString(md, "ws.header", "header"); len(mm) > 0 {
		hd := http.Header{}
		for k, v := range mm {
			hd.Add(k, v)
		}
		l.md.header = hd
	}

	l.md.mptcp = mdutil.GetBool(md, "mptcp")

	l.md.limiterRefreshInterval = mdutil.GetDuration(md, "limiter.refreshInterval")
	if l.md.limiterRefreshInterval == 0 {
		l.md.limiterRefreshInterval = 30 * time.Second
	}
	if l.md.limiterRefreshInterval < time.Second {
		l.md.limiterRefreshInterval = time.Second
	}

	return
}
