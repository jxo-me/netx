package chain

import (
	"github.com/jxo-me/netx/core/auth"
	"github.com/jxo-me/netx/core/bypass"
	"github.com/jxo-me/netx/core/hosts"
	"github.com/jxo-me/netx/core/metadata"
	"github.com/jxo-me/netx/core/resolver"
	"github.com/jxo-me/netx/core/selector"
)

type HTTPNodeSettings struct {
	Host   string
	Header map[string]string
}

type TLSNodeSettings struct {
	ServerName string
	Secure     bool
}

type NodeOptions struct {
	Transport  *Transport
	Bypass     bypass.IBypass
	Resolver   resolver.IResolver
	HostMapper hosts.IHostMapper
	Metadata   metadata.IMetaData
	Host       string
	Protocol   string
	HTTP       *HTTPNodeSettings
	TLS        *TLSNodeSettings
	Auther     auth.IAuthenticator
}

type NodeOption func(*NodeOptions)

func TransportNodeOption(tr *Transport) NodeOption {
	return func(o *NodeOptions) {
		o.Transport = tr
	}
}

func BypassNodeOption(bp bypass.IBypass) NodeOption {
	return func(o *NodeOptions) {
		o.Bypass = bp
	}
}

func ResoloverNodeOption(resolver resolver.IResolver) NodeOption {
	return func(o *NodeOptions) {
		o.Resolver = resolver
	}
}

func HostMapperNodeOption(m hosts.IHostMapper) NodeOption {
	return func(o *NodeOptions) {
		o.HostMapper = m
	}
}

func HostNodeOption(host string) NodeOption {
	return func(o *NodeOptions) {
		o.Host = host
	}
}

func ProtocolNodeOption(protocol string) NodeOption {
	return func(o *NodeOptions) {
		o.Protocol = protocol
	}
}

func MetadataNodeOption(md metadata.IMetaData) NodeOption {
	return func(o *NodeOptions) {
		o.Metadata = md
	}
}

func HTTPNodeOption(httpSettings *HTTPNodeSettings) NodeOption {
	return func(o *NodeOptions) {
		o.HTTP = httpSettings
	}
}

func TLSNodeOption(tlsSettings *TLSNodeSettings) NodeOption {
	return func(o *NodeOptions) {
		o.TLS = tlsSettings
	}
}

func AutherNodeOption(auther auth.IAuthenticator) NodeOption {
	return func(o *NodeOptions) {
		o.Auther = auther
	}
}

type Node struct {
	Name    string
	Addr    string
	marker  selector.IMarker
	options NodeOptions
}

func NewNode(name string, addr string, opts ...NodeOption) *Node {
	var options NodeOptions
	for _, opt := range opts {
		if opt != nil {
			opt(&options)
		}
	}

	return &Node{
		Name:    name,
		Addr:    addr,
		marker:  selector.NewFailMarker(),
		options: options,
	}
}

func (node *Node) Options() *NodeOptions {
	return &node.options
}

// Metadata implements metadadta.IMetaDatable interface.
func (node *Node) Metadata() metadata.IMetaData {
	return node.options.Metadata
}

// Marker implements selector.IMarkable interface.
func (node *Node) Marker() selector.IMarker {
	return node.marker
}

func (node *Node) Copy() *Node {
	n := &Node{}
	*n = *node
	return n
}
