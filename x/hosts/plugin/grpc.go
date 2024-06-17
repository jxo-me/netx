package hosts

import (
	"context"
	"io"
	"net"

	"github.com/jxo-me/netx/core/hosts"
	"github.com/jxo-me/netx/core/logger"
	"github.com/jxo-me/netx/plugin/hosts/proto"
	ctxvalue "github.com/jxo-me/netx/x/ctx"
	"github.com/jxo-me/netx/x/internal/plugin"
	"google.golang.org/grpc"
)

type grpcPlugin struct {
	conn   grpc.ClientConnInterface
	client proto.HostMapperClient
	log    logger.ILogger
}

// NewGRPCPlugin creates a HostMapper plugin based on gRPC.
func NewGRPCPlugin(name string, addr string, opts ...plugin.Option) hosts.IHostMapper {
	var options plugin.Options
	for _, opt := range opts {
		opt(&options)
	}

	log := logger.Default().WithFields(map[string]any{
		"kind":  "hosts",
		"hosts": name,
	})
	conn, err := plugin.NewGRPCConn(addr, &options)
	if err != nil {
		log.Error(err)
	}
	p := &grpcPlugin{
		conn: conn,
		log:  log,
	}
	if conn != nil {
		p.client = proto.NewHostMapperClient(conn)
	}
	return p
}

func (p *grpcPlugin) Lookup(ctx context.Context, network, host string, opts ...hosts.Option) (ips []net.IP, ok bool) {
	p.log.Debugf("lookup %s/%s", host, network)

	if p.client == nil {
		return
	}

	r, err := p.client.Lookup(ctx,
		&proto.LookupRequest{
			Network: network,
			Host:    host,
			Client:  string(ctxvalue.ClientIDFromContext(ctx)),
		})
	if err != nil {
		p.log.Error(err)
		return
	}
	for _, s := range r.Ips {
		if ip := net.ParseIP(s); ip != nil {
			ips = append(ips, ip)
		}
	}
	ok = r.Ok
	return
}

func (p *grpcPlugin) Close() error {
	if p.conn == nil {
		return nil
	}

	if closer, ok := p.conn.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
