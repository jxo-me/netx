package hop

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"github.com/jxo-me/netx/core/chain"
	"github.com/jxo-me/netx/core/hop"
	"github.com/jxo-me/netx/core/logger"
	"github.com/jxo-me/netx/plugin/hop/proto"
	"github.com/jxo-me/netx/x/config"
	node_parser "github.com/jxo-me/netx/x/config/parsing/node"
	ctxvalue "github.com/jxo-me/netx/x/ctx"
	"github.com/jxo-me/netx/x/internal/plugin"
	"google.golang.org/grpc"
)

type grpcPlugin struct {
	name   string
	conn   grpc.ClientConnInterface
	client proto.HopClient
	log    logger.ILogger
}

// NewGRPCPlugin creates a Hop plugin based on gRPC.
func NewGRPCPlugin(name string, addr string, opts ...plugin.Option) hop.IHop {
	var options plugin.Options
	for _, opt := range opts {
		opt(&options)
	}

	log := logger.Default().WithFields(map[string]any{
		"kind": "hop",
		"hop":  name,
	})
	conn, err := plugin.NewGRPCConn(addr, &options)
	if err != nil {
		log.Error(err)
	}

	p := &grpcPlugin{
		name: name,
		conn: conn,
		log:  log,
	}
	if conn != nil {
		p.client = proto.NewHopClient(conn)
	}
	return p
}

func (p *grpcPlugin) Select(ctx context.Context, opts ...hop.SelectOption) *chain.Node {
	if p.client == nil {
		return nil
	}

	var options hop.SelectOptions
	for _, opt := range opts {
		opt(&options)
	}

	r, err := p.client.Select(ctx,
		&proto.SelectRequest{
			Network: options.Network,
			Addr:    options.Addr,
			Host:    options.Host,
			Path:    options.Path,
			Client:  string(ctxvalue.ClientIDFromContext(ctx)),
			Src:     string(ctxvalue.ClientAddrFromContext(ctx)),
		})
	if err != nil {
		p.log.Error(err)
		return nil
	}

	if r.Node == nil {
		return nil
	}

	var cfg config.NodeConfig
	if err := json.NewDecoder(bytes.NewReader(r.Node)).Decode(&cfg); err != nil {
		p.log.Error(err)
		return nil
	}

	node, err := node_parser.ParseNode(p.name, &cfg, logger.Default())
	if err != nil {
		p.log.Error(err)
		return nil
	}
	return node
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
