package handler

import (
	"crypto/tls"
	"net/url"

	"github.com/jxo-me/netx/core/auth"
	"github.com/jxo-me/netx/core/bypass"
	"github.com/jxo-me/netx/core/chain"
	"github.com/jxo-me/netx/core/limiter/rate"
	"github.com/jxo-me/netx/core/logger"
	"github.com/jxo-me/netx/core/metadata"
)

type Options struct {
	Bypass      bypass.Bypass
	Router      *chain.Router
	Auth        *url.Userinfo
	Auther      auth.Authenticator
	RateLimiter rate.IRateLimiter
	TLSConfig   *tls.Config
	Logger      logger.ILogger
	Service     string
}

type Option func(opts *Options)

func BypassOption(bypass bypass.Bypass) Option {
	return func(opts *Options) {
		opts.Bypass = bypass
	}
}

func RouterOption(router *chain.Router) Option {
	return func(opts *Options) {
		opts.Router = router
	}
}

func AuthOption(auth *url.Userinfo) Option {
	return func(opts *Options) {
		opts.Auth = auth
	}
}

func AutherOption(auther auth.Authenticator) Option {
	return func(opts *Options) {
		opts.Auther = auther
	}
}

func RateLimiterOption(limiter rate.IRateLimiter) Option {
	return func(opts *Options) {
		opts.RateLimiter = limiter
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

func ServiceOption(service string) Option {
	return func(opts *Options) {
		opts.Service = service
	}
}

type HandleOptions struct {
	Metadata metadata.IMetaData
}

type HandleOption func(opts *HandleOptions)

func MetadataHandleOption(md metadata.IMetaData) HandleOption {
	return func(opts *HandleOptions) {
		opts.Metadata = md
	}
}
