package relay

import (
	"time"

	"github.com/google/uuid"
	mdata "github.com/jxo-me/netx/sdk/core/metadata"
	mdutil "github.com/jxo-me/netx/sdk/core/metadata/util"
	"github.com/jxo-me/netx/sdk/relay"
)

type metadata struct {
	connectTimeout time.Duration
	noDelay        bool
	tunnelID       relay.TunnelID
}

func (c *relayConnector) parseMetadata(md mdata.IMetaData) (err error) {
	const (
		connectTimeout = "connectTimeout"
		noDelay        = "nodelay"
	)

	c.md.connectTimeout = mdutil.GetDuration(md, connectTimeout)
	c.md.noDelay = mdutil.GetBool(md, noDelay)

	if s := mdutil.GetString(md, "tunnelID", "tunnel.id"); s != "" {
		uuid, err := uuid.Parse(s)
		if err != nil {
			return err
		}
		c.md.tunnelID = relay.NewTunnelID(uuid[:])
	}

	return
}
