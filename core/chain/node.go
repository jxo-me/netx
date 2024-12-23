package chain

import (
	"regexp"

	"github.com/jxo-me/netx/core/auth"
	"github.com/jxo-me/netx/core/bypass"
	"github.com/jxo-me/netx/core/hosts"
	"github.com/jxo-me/netx/core/metadata"
	"github.com/jxo-me/netx/core/resolver"
	"github.com/jxo-me/netx/core/routing"
	"github.com/jxo-me/netx/core/selector"
)

type NodeFilterSettings struct {
	Protocol string
	Host     string
	Path     string
}

type HTTPURLRewriteSetting struct {
	Pattern     *regexp.Regexp
	Replacement string
}

type HTTPBodyRewriteSettings struct {
	Type        string
	Pattern     *regexp.Regexp
	Replacement []byte
}

type HTTPNodeSettings struct {
	Host                string
	RequestHeader       map[string]string
	ResponseHeader      map[string]string
	Auther              auth.IAuthenticator
	RewriteURL          []HTTPURLRewriteSetting
	RewriteResponseBody []HTTPBodyRewriteSettings
}

type TLSNodeSettings struct {
	ServerName string
	Secure     bool
	Options    struct {
		MinVersion   string
		MaxVersion   string
		CipherSuites []string
		ALPN         []string
	}
}

type NodeOptions struct {
	Network    string
	Transport  Transporter
	Bypass     bypass.IBypass
	Resolver   resolver.IResolver
	HostMapper hosts.IHostMapper
	Filter     *NodeFilterSettings
	HTTP       *HTTPNodeSettings
	TLS        *TLSNodeSettings
	Metadata   metadata.IMetaData
	Matcher    routing.Matcher
	Priority   int
}

type NodeOption func(*NodeOptions)

func TransportNodeOption(tr Transporter) NodeOption {
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

func NetworkNodeOption(network string) NodeOption {
	return func(o *NodeOptions) {
		o.Network = network
	}
}

func NodeFilterOption(filter *NodeFilterSettings) NodeOption {
	return func(o *NodeOptions) {
		o.Filter = filter
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

func MetadataNodeOption(md metadata.IMetaData) NodeOption {
	return func(o *NodeOptions) {
		o.Metadata = md
	}
}

func MatcherNodeOption(matcher routing.Matcher) NodeOption {
	return func(o *NodeOptions) {
		o.Matcher = matcher
	}
}

func PriorityNodeOption(priority int) NodeOption {
	return func(o *NodeOptions) {
		o.Priority = priority
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
