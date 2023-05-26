package auth

import (
	"context"

	xlogger "github.com/jxo-me/netx/sdk/core/logger"
	"github.com/jxo-me/netx/sdk/plugin/auth/proto"
)

type pluginAuthenticator struct {
	client  proto.AuthenticatorClient
	options options
}

// NewPluginAuthenticator creates an Authenticator that authenticates client by plugin.
func NewPluginAuthenticator(opts ...Option) IAuthenticator {
	var options options
	for _, opt := range opts {
		opt(&options)
	}
	if options.logger == nil {
		options.logger = xlogger.Nop()
	}

	p := &pluginAuthenticator{
		options: options,
	}
	if options.client != nil {
		p.client = proto.NewAuthenticatorClient(options.client)
	}
	return p
}

// Authenticate checks the validity of the provided user-password pair.
func (p *pluginAuthenticator) Authenticate(ctx context.Context, user, password string) bool {
	if p.client == nil {
		return false
	}

	r, err := p.client.Authenticate(ctx,
		&proto.AuthenticateRequest{
			Username: user,
			Password: password,
		})
	if err != nil {
		p.options.logger.Error(err)
		return false
	}
	return r.Ok
}

func (p *pluginAuthenticator) Close() error {
	if p.options.client != nil {
		return p.options.client.Close()
	}
	return nil
}
