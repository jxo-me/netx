package ingress

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/jxo-me/netx/core/ingress"
	"github.com/jxo-me/netx/core/logger"
	"github.com/jxo-me/netx/plugin/ingress/proto"
	"github.com/jxo-me/netx/x/internal/plugin"
	"google.golang.org/grpc"
)

type grpcPlugin struct {
	conn   grpc.ClientConnInterface
	client proto.IngressClient
	log    logger.ILogger
}

// NewGRPCPlugin creates an Ingress plugin based on gRPC.
func NewGRPCPlugin(name string, addr string, opts ...plugin.Option) ingress.IIngress {
	var options plugin.Options
	for _, opt := range opts {
		opt(&options)
	}

	log := logger.Default().WithFields(map[string]any{
		"kind":    "ingress",
		"ingress": name,
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
		p.client = proto.NewIngressClient(conn)
	}
	return p
}

func (p *grpcPlugin) GetRule(ctx context.Context, host string, opts ...ingress.Option) *ingress.Rule {
	if p.client == nil {
		return nil
	}

	r, err := p.client.GetRule(ctx,
		&proto.GetRuleRequest{
			Host: host,
		})
	if err != nil {
		p.log.Error(err)
		return nil
	}
	if r.Endpoint == "" {
		return nil
	}
	return &ingress.Rule{
		Hostname: host,
		Endpoint: r.Endpoint,
	}
}

func (p *grpcPlugin) SetRule(ctx context.Context, rule *ingress.Rule, opts ...ingress.Option) bool {
	if p.client == nil || rule == nil {
		return false
	}

	r, _ := p.client.SetRule(ctx, &proto.SetRuleRequest{
		Host:     rule.Hostname,
		Endpoint: rule.Endpoint,
	})
	if r == nil {
		return false
	}

	return r.Ok
}

func (p *grpcPlugin) Close() error {
	if closer, ok := p.conn.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

type httpPluginGetRuleRequest struct {
	Host string `json:"host"`
}

type httpPluginGetRuleResponse struct {
	Endpoint string `json:"endpoint"`
}

type httpPluginSetRuleRequest struct {
	Host     string `json:"host"`
	Endpoint string `json:"endpoint"`
}

type httpPluginSetRuleResponse struct {
	OK bool `json:"ok"`
}

type httpPlugin struct {
	url    string
	client *http.Client
	header http.Header
	log    logger.ILogger
}

// NewHTTPPlugin creates an Ingress plugin based on HTTP.
func NewHTTPPlugin(name string, url string, opts ...plugin.Option) ingress.IIngress {
	var options plugin.Options
	for _, opt := range opts {
		opt(&options)
	}

	return &httpPlugin{
		url:    url,
		client: plugin.NewHTTPClient(&options),
		header: options.Header,
		log: logger.Default().WithFields(map[string]any{
			"kind":    "ingress",
			"ingress": name,
		}),
	}
}

func (p *httpPlugin) GetRule(ctx context.Context, host string, opts ...ingress.Option) *ingress.Rule {
	if p.client == nil {
		return nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.url, nil)
	if err != nil {
		return nil
	}
	if p.header != nil {
		req.Header = p.header.Clone()
	}
	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()
	q.Set("host", host)
	req.URL.RawQuery = q.Encode()

	resp, err := p.client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	res := httpPluginGetRuleResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil
	}
	if res.Endpoint == "" {
		return nil
	}
	return &ingress.Rule{
		Hostname: host,
		Endpoint: res.Endpoint,
	}
}

func (p *httpPlugin) SetRule(ctx context.Context, rule *ingress.Rule, opts ...ingress.Option) bool {
	if p.client == nil || rule == nil {
		return false
	}

	rb := httpPluginSetRuleRequest{
		Host:     rule.Hostname,
		Endpoint: rule.Endpoint,
	}
	v, err := json.Marshal(&rb)
	if err != nil {
		return false
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, p.url, bytes.NewReader(v))
	if err != nil {
		return false
	}

	if p.header != nil {
		req.Header = p.header.Clone()
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	res := httpPluginSetRuleResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return false
	}
	return res.OK
}
