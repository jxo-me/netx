package auth

import "context"

type Options struct{}
type Option func(opts *Options)
// IAuthenticator is an interface for user authentication.
type IAuthenticator interface {
	Authenticate(ctx context.Context, user, password string, opts ...Option) (id string, ok bool)
}

type authenticatorGroup struct {
	authers []IAuthenticator
}

func AuthenticatorGroup(authers ...IAuthenticator) IAuthenticator {
	return &authenticatorGroup{
		authers: authers,
	}
}

func (p *authenticatorGroup) Authenticate(ctx context.Context, user, password string, opts ...Option) (string, bool) {
	if len(p.authers) == 0 {
		return "", false
	}
	for _, auther := range p.authers {
		if auther == nil {
			continue
		}

		if id, ok := auther.Authenticate(ctx, user, password, opts...); ok {
			return id, ok
		}
	}
	return "", false
}
