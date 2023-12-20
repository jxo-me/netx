package file

import (
	mdata "github.com/jxo-me/netx/core/metadata"
	mdutil "github.com/jxo-me/netx/core/metadata/util"
)

type metadata struct {
	dir string
}

func (h *fileHandler) parseMetadata(md mdata.IMetaData) (err error) {
	h.md.dir = mdutil.GetString(md, "file.dir", "dir")
	return
}
