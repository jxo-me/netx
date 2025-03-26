package tunnel

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"net"

	"github.com/google/uuid"
	"github.com/jxo-me/netx/core/ingress"
	"github.com/jxo-me/netx/core/logger"
	"github.com/jxo-me/netx/core/observer/stats"
	"github.com/jxo-me/netx/core/sd"
	"github.com/jxo-me/netx/relay"
	"github.com/jxo-me/netx/x/internal/util/mux"
)

func (h *tunnelHandler) handleBind(ctx context.Context, conn net.Conn, network, address string, tunnelID relay.TunnelID, log logger.ILogger) (err error) {
	resp := relay.Response{
		Version: relay.Version1,
		Status:  relay.StatusOK,
	}

	uuid, err := uuid.NewRandom()
	if err != nil {
		resp.Status = relay.StatusInternalServerError
		resp.WriteTo(conn)
		return
	}
	connectorID := relay.NewConnectorID(uuid[:])
	if network == "udp" {
		connectorID = relay.NewUDPConnectorID(uuid[:])
	}
	// copy weight from tunnelID
	connectorID = connectorID.SetWeight(tunnelID.Weight())

	v := md5.Sum([]byte(tunnelID.String()))
	endpoint := hex.EncodeToString(v[:8])

	host, port, _ := net.SplitHostPort(address)
	if host == "" || h.md.ingress == nil {
		host = endpoint
	} else if host != endpoint {
		if rule := h.md.ingress.GetRule(ctx, host); rule != nil && rule.Endpoint != tunnelID.String() {
			host = endpoint
		}
	}
	addr := net.JoinHostPort(host, port)

	af := &relay.AddrFeature{}
	err = af.ParseFrom(addr)
	if err != nil {
		log.Warn(err)
	}
	resp.Features = append(resp.Features, af,
		&relay.TunnelFeature{
			ID: connectorID,
		},
	)
	resp.WriteTo(conn)

	// Upgrade connection to multiplex session.
	session, err := mux.ClientSession(conn, h.md.muxCfg)
	if err != nil {
		return
	}

	var stats stats.Stats
	if h.stats != nil {
		stats = h.stats.Stats(tunnelID.String())
	}

	c := NewConnector(connectorID, tunnelID, h.id, session, &ConnectorOptions{
		service: h.options.Service,
		sd:      h.md.sd,
		stats:   stats,
		limiter: h.limiter,
	})

	h.pool.Add(tunnelID, c, h.md.tunnelTTL)
	if h.md.ingress != nil {
		h.md.ingress.SetRule(ctx, &ingress.Rule{
			Hostname: endpoint,
			Endpoint: tunnelID.String(),
		})
		if host != "" {
			h.md.ingress.SetRule(ctx, &ingress.Rule{
				Hostname: host,
				Endpoint: tunnelID.String(),
			})
		}
	}
	if h.md.sd != nil {
		err := h.md.sd.Register(ctx, &sd.Service{
			ID:      connectorID.String(),
			Name:    tunnelID.String(),
			Node:    h.id,
			Network: network,
			Address: h.md.entryPoint,
		})
		if err != nil {
			h.log.Error(err)
		}
	}

	log.Debugf("%s/%s: tunnel=%s, connector=%s, weight=%d established", addr, network, tunnelID, connectorID, connectorID.Weight())

	return
}
