package recorder

import "context"

type IRecorder interface {
	Record(ctx context.Context, b []byte) error
}

type RecorderObject struct {
	Recorder IRecorder
	Record   string
}

const (
	RecorderServiceClientAddress          = "recorder.service.client.address"
	RecorderServiceRouterDialAddress      = "recorder.service.router.dial.address"
	RecorderServiceRouterDialAddressError = "recorder.service.router.dial.address.error"
)
