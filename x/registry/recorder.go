package registry

import (
	"context"

	"github.com/jxo-me/netx/core/recorder"
)

type RecorderRegistry struct {
	registry[recorder.IRecorder]
}

func (r *RecorderRegistry) Register(name string, v recorder.IRecorder) error {
	return r.registry.Register(name, v)
}

func (r *RecorderRegistry) Get(name string) recorder.IRecorder {
	if name != "" {
		return &recorderWrapper{name: name, r: r}
	}
	return nil
}

func (r *RecorderRegistry) get(name string) recorder.IRecorder {
	return r.registry.Get(name)
}

type recorderWrapper struct {
	name string
	r    *RecorderRegistry
}

func (w *recorderWrapper) Record(ctx context.Context, b []byte, opts ...recorder.RecordOption) error {
	v := w.r.get(w.name)
	if v == nil {
		return nil
	}
	return v.Record(ctx, b, opts...)
}
