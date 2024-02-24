package observer

import "context"

type Options struct{}

type Option func(opts *Options)

type IObserver interface {
	Observe(ctx context.Context, events []Event, opts ...Option) error
}

type EventType string

const (
	EventStatus EventType = "status"
	EventStats  EventType = "stats"
)

type Event interface {
	Type() EventType
}
