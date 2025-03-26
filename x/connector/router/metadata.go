package router

import (
	"errors"
	"time"

	"github.com/google/uuid"
	mdata "github.com/jxo-me/netx/core/metadata"
	"github.com/jxo-me/netx/relay"
	mdutil "github.com/jxo-me/netx/x/metadata/util"
)

var (
	ErrInvalidRouterID = errors.New("router: invalid router ID")
)

type metadata struct {
	connectTimeout time.Duration
	routerID       relay.TunnelID
}

func (c *routerConnector) parseMetadata(md mdata.IMetaData) (err error) {
	c.md.connectTimeout = mdutil.GetDuration(md, "connectTimeout")

	if s := mdutil.GetString(md, "router.id"); s != "" {
		uuid, err := uuid.Parse(s)
		if err != nil {
			return err
		}
		c.md.routerID = relay.NewTunnelID(uuid[:])
	}

	return
}
