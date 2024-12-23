package recorder

import (
	"context"
)

type RecordOptions struct {
	Metadata any
}

type RecordOption func(opts *RecordOptions)

func MetadataRecordOption(md any) RecordOption {
	return func(opts *RecordOptions) {
		opts.Metadata = md
	}
}

type IRecorder interface {
	Record(ctx context.Context, b []byte, opts ...RecordOption) error
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
	HTTPBody        bool
	MaxBodySize     int
}

const (
	RecorderServiceClientAddress          = "recorder.service.client.address"
	RecorderServiceRouterDialAddress      = "recorder.service.router.dial.address"
	RecorderServiceRouterDialAddressError = "recorder.service.router.dial.address.error"
)
