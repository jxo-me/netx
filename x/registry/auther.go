package registry

import (
	"context"

	"github.com/jxo-me/netx/core/auth"
)

type AutherRegistry struct {
	registry[auth.IAuthenticator]
}

func (r *AutherRegistry) Register(name string, v auth.IAuthenticator) error {
	return r.registry.Register(name, v)
}

func (r *AutherRegistry) Get(name string) auth.IAuthenticator {
	if name != "" {
		return &autherWrapper{name: name, r: r}
	}
	return nil
}

func (r *AutherRegistry) get(name string) auth.IAuthenticator {
	return r.registry.Get(name)
}

type autherWrapper struct {
	name string
	r    *AutherRegistry
}

func (w *autherWrapper) Authenticate(ctx context.Context, user, password string) (string, bool) {
	v := w.r.get(w.name)
	if v == nil {
		return "", true
	}
	return v.Authenticate(ctx, user, password)
}
