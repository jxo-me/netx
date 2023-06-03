package api

import "github.com/jxo-me/netx/core/auth"

type options struct {
	accessLog  bool
	pathPrefix string
	auther     auth.IAuthenticator
	botEnable  bool
	domain     string
	token      string
}

type Option func(*options)

func PathPrefixOption(pathPrefix string) Option {
	return func(o *options) {
		o.pathPrefix = pathPrefix
	}
}

func AccessLogOption(enable bool) Option {
	return func(o *options) {
		o.accessLog = enable
	}
}

func AutherOption(auther auth.IAuthenticator) Option {
	return func(o *options) {
		o.auther = auther
	}
}

func DomainOption(domain string) Option {
	return func(o *options) {
		o.domain = domain
	}
}

func BotEnableOption(enable bool) Option {
	return func(o *options) {
		o.botEnable = enable
	}
}

func TokenOption(token string) Option {
	return func(o *options) {
		o.token = token
	}
}
