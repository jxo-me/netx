package tunnel

import (
	"errors"
	"time"

	"github.com/google/uuid"
	mdata "github.com/jxo-me/netx/core/metadata"
	"github.com/jxo-me/netx/relay"
	"github.com/jxo-me/netx/x/internal/util/mux"
	mdutil "github.com/jxo-me/netx/x/metadata/util"
)

var (
	ErrInvalidTunnelID = errors.New("tunnel: invalid tunnel ID")
)

type metadata struct {
	connectTimeout time.Duration
	tunnelID       relay.TunnelID
	muxCfg         *mux.Config
}

func (c *tunnelConnector) parseMetadata(md mdata.IMetaData) (err error) {
	c.md.connectTimeout = mdutil.GetDuration(md, "connectTimeout")

	if s := mdutil.GetString(md, "tunnelID", "tunnel.id"); s != "" {
		uuid, err := uuid.Parse(s)
		if err != nil {
			return err
		}
		c.md.tunnelID = relay.NewTunnelID(uuid[:])
	}

	if c.md.tunnelID.IsZero() {
		uuid, err := uuid.NewUUID()
		if err != nil {
			return err
		}
		c.md.tunnelID = relay.NewTunnelID(uuid[:])
	}

	if weight := mdutil.GetInt(md, "tunnel.weight"); weight > 0 {
		c.md.tunnelID = c.md.tunnelID.SetWeight(uint8(weight))
	}

	c.md.muxCfg = &mux.Config{
		Version:           mdutil.GetInt(md, "mux.version"),
		KeepAliveInterval: mdutil.GetDuration(md, "mux.keepaliveInterval"),
		KeepAliveDisabled: mdutil.GetBool(md, "mux.keepaliveDisabled"),
		KeepAliveTimeout:  mdutil.GetDuration(md, "mux.keepaliveTimeout"),
		MaxFrameSize:      mdutil.GetInt(md, "mux.maxFrameSize"),
		MaxReceiveBuffer:  mdutil.GetInt(md, "mux.maxReceiveBuffer"),
		MaxStreamBuffer:   mdutil.GetInt(md, "mux.maxStreamBuffer"),
	}
	if c.md.muxCfg.Version == 0 {
		c.md.muxCfg.Version = 2
	}
	if c.md.muxCfg.MaxStreamBuffer == 0 {
		c.md.muxCfg.MaxStreamBuffer = 1048576
	}

	return
}
