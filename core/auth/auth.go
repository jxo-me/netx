package auth

import "context"

type Options struct{}
type Option func(opts *Options)
// IAuthenticator is an interface for user authentication.
type IAuthenticator interface {
	Authenticate(ctx context.Context, user, password string, opts ...Option) (id string, ok bool)
}
