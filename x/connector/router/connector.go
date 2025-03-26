package router

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/jxo-me/netx/core/connector"
	md "github.com/jxo-me/netx/core/metadata"
	"github.com/jxo-me/netx/relay"
	ctxvalue "github.com/jxo-me/netx/x/ctx"
)

type routerConnector struct {
	md      metadata
	options connector.Options
}

func NewConnector(opts ...connector.Option) connector.IConnector {
	options := connector.Options{}
	for _, opt := range opts {
		opt(&options)
	}

	return &routerConnector{
		options: options,
	}
}

func (c *routerConnector) Init(md md.IMetaData) (err error) {
	return c.parseMetadata(md)
}

func (c *routerConnector) Connect(ctx context.Context, conn net.Conn, network, address string, opts ...connector.ConnectOption) (net.Conn, error) {
	log := c.options.Logger.WithFields(map[string]any{
		"remote":  conn.RemoteAddr().String(),
		"local":   conn.LocalAddr().String(),
		"network": network,
		"address": address,
		"sid":     string(ctxvalue.SidFromContext(ctx)),
	})
	log.Debugf("connect %s/%s", address, network)

	if c.md.connectTimeout > 0 {
		conn.SetDeadline(time.Now().Add(c.md.connectTimeout))
		defer conn.SetDeadline(time.Time{})
	}

	req := relay.Request{
		Version: relay.Version1,
		Cmd:     relay.CmdAssociate,
	}

	if c.options.Auth != nil {
		pwd, _ := c.options.Auth.Password()
		req.Features = append(req.Features, &relay.UserAuthFeature{
			Username: c.options.Auth.Username(),
			Password: pwd,
		})
	}

	switch network {
	case "udp":
		req.Features = append(req.Features, &relay.NetworkFeature{
			Network: relay.NetworkUDP,
		})

	case "ip":
		req.Features = append(req.Features, &relay.NetworkFeature{
			Network: relay.NetworkIP,
		})
	}

	srcAddr := conn.LocalAddr().String()
	if v := ctxvalue.ClientAddrFromContext(ctx); v != "" {
		srcAddr = string(v)
	}

	af := &relay.AddrFeature{}
	af.ParseFrom(srcAddr)
	req.Features = append(req.Features, af) // src address

	af = &relay.AddrFeature{}
	af.ParseFrom(address)
	req.Features = append(req.Features, af) // dst address

	if !c.md.routerID.IsZero() {
		req.Features = append(req.Features, &relay.TunnelFeature{
			ID: c.md.routerID,
		})
	}

	if _, err := req.WriteTo(conn); err != nil {
		return nil, err
	}
	// drain the response
	if err := readResponse(conn); err != nil {
		return nil, err
	}

	switch network {
	case "udp", "ip":
		conn = &packetConn{
			Conn: conn,
		}
	default:
		err := fmt.Errorf("network %s is unsupported", network)
		log.Error(err)
		return nil, err
	}

	log.Debugf("connect to router %s@%s OK", address, c.md.routerID)
	return conn, nil
}
