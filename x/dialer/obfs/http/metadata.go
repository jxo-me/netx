package http

import (
	"net/http"

	mdata "github.com/jxo-me/netx/core/metadata"
	mdutil "github.com/jxo-me/netx/x/metadata/util"
)

const (
	defaultPath = "/"
)

type metadata struct {
	host   string
	path   string
	header http.Header
}

func (d *obfsHTTPDialer) parseMetadata(md mdata.IMetaData) (err error) {
	d.md.host = mdutil.GetString(md, "obfs.host", "host")
	d.md.path = mdutil.GetString(md, "obfs.path", "path")
	if d.md.path == "" {
		d.md.path = defaultPath
	}

	if m := mdutil.GetStringMapString(md, "obfs.header", "header"); len(m) > 0 {
		h := http.Header{}
		for k, v := range m {
			h.Add(k, v)
		}
		d.md.header = h
	}
	return
}
