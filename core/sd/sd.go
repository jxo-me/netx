package sd

import (
	"context"
)

type Options struct{}

type Option func(opts *Options)

type Service struct {
	ID      string
	Name    string
	Node    string
	Network string
	Address string
}

// ISD is the service discovery interface.
type ISD interface {
	Register(ctx context.Context, service *Service, opts ...Option) error
	Deregister(ctx context.Context, service *Service) error
	Renew(ctx context.Context, service *Service) error
	Get(ctx context.Context, name string) ([]*Service, error)
}
