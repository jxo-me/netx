package relay

import (
	"fmt"
	"net"
	"strconv"

	"github.com/jxo-me/netx/sdk/core/internal/util/mux"
	"github.com/jxo-me/netx/sdk/core/logger"
	mdata "github.com/jxo-me/netx/sdk/core/metadata"
	mdx "github.com/jxo-me/netx/sdk/core/metadata"
	"github.com/jxo-me/netx/sdk/relay"
)

type bindListener struct {
	network string
	addr    net.Addr
	session *mux.Session
	logger  logger.ILogger
}

func (p *bindListener) Accept() (net.Conn, error) {
	cc, err := p.session.Accept()
	if err != nil {
		return nil, err
	}

	conn, err := p.getPeerConn(cc)
	if err != nil {
		cc.Close()
		p.logger.Errorf("get peer failed: %s", err)
		return nil, err
	}

	return conn, nil
}

func (p *bindListener) getPeerConn(conn net.Conn) (net.Conn, error) {
	// second reply, peer connected
	resp := relay.Response{}
	if _, err := resp.ReadFrom(conn); err != nil {
		return nil, err
	}

	if resp.Status != relay.StatusOK {
		err := fmt.Errorf("peer connect failed")
		return nil, err
	}

	var address, host string
	for _, f := range resp.Features {
		if f.Type() == relay.FeatureAddr {
			if fa, ok := f.(*relay.AddrFeature); ok {
				v := net.JoinHostPort(fa.Host, strconv.Itoa(int(fa.Port)))
				if address != "" {
					host = v
				} else {
					address = v
				}
			}
		}
	}

	raddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}

	var md mdata.IMetaData
	if host != "" {
		md = mdx.NewMetadata(map[string]any{"host": host})
	}

	if p.network == "udp" {
		return &bindUDPConn{
			Conn:       conn,
			localAddr:  p.addr,
			remoteAddr: raddr,
			md:         md,
		}, nil
	}

	cn := &bindConn{
		Conn:       conn,
		localAddr:  p.addr,
		remoteAddr: raddr,
		md:         md,
	}
	return cn, nil
}

func (p *bindListener) Addr() net.Addr {
	return p.addr
}

func (p *bindListener) Close() error {
	return p.session.Close()
}
