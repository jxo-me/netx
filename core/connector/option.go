package connector

import (
	"crypto/tls"
	"net/url"
	"time"

	xnet "github.com/jxo-me/netx/core/common/net"
	"github.com/jxo-me/netx/core/logger"
)

type Options struct {
	Auth      *url.Userinfo
	TLSConfig *tls.Config
	Logger    logger.ILogger
}

type Option func(opts *Options)

func AuthOption(auth *url.Userinfo) Option {
	return func(opts *Options) {
		opts.Auth = auth
	}
}

func TLSConfigOption(tlsConfig *tls.Config) Option {
	return func(opts *Options) {
		opts.TLSConfig = tlsConfig
	}
}

func LoggerOption(logger logger.ILogger) Option {
	return func(opts *Options) {
		opts.Logger = logger
	}
}

type ConnectOptions struct {
	Dialer xnet.Dialer
}

type ConnectOption func(opts *ConnectOptions)

func DialerConnectOption(dialer xnet.Dialer) ConnectOption {
	return func(opts *ConnectOptions) {
		opts.Dialer = dialer
	}
}

type BindOptions struct {
	Mux               bool
	Backlog           int
	UDPDataQueueSize  int
	UDPDataBufferSize int
	UDPConnTTL        time.Duration
}

type BindOption func(opts *BindOptions)

func MuxBindOption(mux bool) BindOption {
	return func(opts *BindOptions) {
		opts.Mux = mux
	}
}

func BacklogBindOption(backlog int) BindOption {
	return func(opts *BindOptions) {
		opts.Backlog = backlog
	}
}

func UDPDataQueueSizeBindOption(size int) BindOption {
	return func(opts *BindOptions) {
		opts.UDPDataQueueSize = size
	}
}

func UDPDataBufferSizeBindOption(size int) BindOption {
	return func(opts *BindOptions) {
		opts.UDPDataBufferSize = size
	}
}

func UDPConnTTLBindOption(ttl time.Duration) BindOption {
	return func(opts *BindOptions) {
		opts.UDPConnTTL = ttl
	}
}
