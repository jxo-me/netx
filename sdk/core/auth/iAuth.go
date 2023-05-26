package auth

import "context"

// IAuthenticator is an interface for user authentication.
type IAuthenticator interface {
	Authenticate(ctx context.Context, user, password string) bool
}

type authenticatorGroup struct {
	authers []IAuthenticator
}

func AuthenticatorGroup(authers ...IAuthenticator) IAuthenticator {
	return &authenticatorGroup{
		authers: authers,
	}
}

func (p *authenticatorGroup) Authenticate(ctx context.Context, user, password string) bool {
	if len(p.authers) == 0 {
		return true
	}
	for _, auther := range p.authers {
		if auther != nil && auther.Authenticate(ctx, user, password) {
			return true
		}
	}
	return false
}
