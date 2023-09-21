package recorder

import (
	"context"
)

type IRecorder interface {
	Record(ctx context.Context, b []byte) error
}

type RecorderObject struct {
	Recorder IRecorder
	Record   string
	Options  *Options
}

type Options struct {
	Direction       bool
	TimestampFormat string
	Hexdump         bool
}

const (
	RecorderServiceClientAddress          = "recorder.service.client.address"
	RecorderServiceRouterDialAddress      = "recorder.service.router.dial.address"
	RecorderServiceRouterDialAddressError = "recorder.service.router.dial.address.error"
)
