package kcp

import (
	"encoding/json"
	"time"

	kcp_util "github.com/jxo-me/netx/sdk/core/internal/util/kcp"
	mdata "github.com/jxo-me/netx/sdk/core/metadata"
	mdutil "github.com/jxo-me/netx/sdk/core/metadata/util"
)

type metadata struct {
	handshakeTimeout time.Duration
	config           *kcp_util.Config
}

func (d *kcpDialer) parseMetadata(md mdata.IMetaData) (err error) {
	const (
		config           = "config"
		configFile       = "c"
		handshakeTimeout = "handshakeTimeout"
	)

	if file := mdutil.GetString(md, configFile); file != "" {
		d.md.config, err = kcp_util.ParseFromFile(file)
		if err != nil {
			return
		}
	}

	if m := mdutil.GetStringMap(md, config); len(m) > 0 {
		b, err := json.Marshal(m)
		if err != nil {
			return err
		}
		cfg := &kcp_util.Config{}
		if err := json.Unmarshal(b, cfg); err != nil {
			return err
		}
		d.md.config = cfg
	}
	if d.md.config == nil {
		d.md.config = kcp_util.DefaultConfig
	}

	d.md.handshakeTimeout = mdutil.GetDuration(md, handshakeTimeout)
	return
}
