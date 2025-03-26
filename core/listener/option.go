package listener

import (
	"crypto/tls"
	"net/url"

	"github.com/jxo-me/netx/core/admission"
	"github.com/jxo-me/netx/core/auth"
	"github.com/jxo-me/netx/core/chain"
	"github.com/jxo-me/netx/core/limiter/conn"
	"github.com/jxo-me/netx/core/limiter/traffic"
	"github.com/jxo-me/netx/core/logger"
	"github.com/jxo-me/netx/core/observer/stats"
)

type Options struct {
	Addr           string
	Auther         auth.IAuthenticator
	Auth           *url.Userinfo
	TLSConfig      *tls.Config
	Admission      admission.IAdmission
	TrafficLimiter traffic.ITrafficLimiter
	ConnLimiter    conn.IConnLimiter
	Chain          chain.IChainer
	Stats          stats.Stats
	Logger         logger.ILogger
	Service        string
	ProxyProtocol  int
	Netns          string
	Router         chain.Router
}

type Option func(opts *Options)

func AddrOption(addr string) Option {
	return func(opts *Options) {
		opts.Addr = addr
	}
}

func AutherOption(auther auth.IAuthenticator) Option {
	return func(opts *Options) {
		opts.Auther = auther
	}
}

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

func AdmissionOption(admission admission.IAdmission) Option {
	return func(opts *Options) {
		opts.Admission = admission
	}
}

func TrafficLimiterOption(limiter traffic.ITrafficLimiter) Option {
	return func(opts *Options) {
		opts.TrafficLimiter = limiter
	}
}

func ConnLimiterOption(limiter conn.IConnLimiter) Option {
	return func(opts *Options) {
		opts.ConnLimiter = limiter
	}
}

func StatsOption(stats stats.Stats) Option {
	return func(opts *Options) {
		opts.Stats = stats
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

func ProxyProtocolOption(ppv int) Option {
	return func(opts *Options) {
		opts.ProxyProtocol = ppv
	}
}

func NetnsOption(netns string) Option {
	return func(opts *Options) {
		opts.Netns = netns
	}
}

func RouterOption(router chain.Router) Option {
	return func(opts *Options) {
		opts.Router = router
	}
}
