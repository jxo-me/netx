package handler

import (
	"crypto/tls"
	"net/url"

	"github.com/jxo-me/netx/core/auth"
	"github.com/jxo-me/netx/core/bypass"
	"github.com/jxo-me/netx/core/chain"
	"github.com/jxo-me/netx/core/limiter/rate"
	"github.com/jxo-me/netx/core/limiter/traffic"
	"github.com/jxo-me/netx/core/logger"
	"github.com/jxo-me/netx/core/metadata"
	"github.com/jxo-me/netx/core/observer"
	"github.com/jxo-me/netx/core/recorder"
)

type Options struct {
	Bypass      bypass.IBypass
	Router      chain.Router
	Auth        *url.Userinfo
	Auther      auth.IAuthenticator
	RateLimiter rate.IRateLimiter
	Limiter     traffic.ITrafficLimiter
	TLSConfig   *tls.Config
	Logger      logger.ILogger
	Observer    observer.IObserver
	Recorders   []recorder.RecorderObject
	Service     string
	Netns       string
}

type Option func(opts *Options)

func BypassOption(bypass bypass.IBypass) Option {
	return func(opts *Options) {
		opts.Bypass = bypass
	}
}

func RouterOption(router chain.Router) Option {
	return func(opts *Options) {
		opts.Router = router
	}
}

func AuthOption(auth *url.Userinfo) Option {
	return func(opts *Options) {
		opts.Auth = auth
	}
}

func AutherOption(auther auth.IAuthenticator) Option {
	return func(opts *Options) {
		opts.Auther = auther
	}
}

func RateLimiterOption(limiter rate.IRateLimiter) Option {
	return func(opts *Options) {
		opts.RateLimiter = limiter
	}
}

func TrafficLimiterOption(limiter traffic.ITrafficLimiter) Option {
	return func(opts *Options) {
		opts.Limiter = limiter
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

func ObserverOption(observer observer.IObserver) Option {
	return func(opts *Options) {
		opts.Observer = observer
	}
}

func RecordersOption(recorders ...recorder.RecorderObject) Option {
	return func(o *Options) {
		o.Recorders = recorders
	}
}

func ServiceOption(service string) Option {
	return func(opts *Options) {
		opts.Service = service
	}
}

func NetnsOption(netns string) Option {
	return func(opts *Options) {
		opts.Netns = netns
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
