package registry

import (
	"context"

	"github.com/jxo-me/netx/core/recorder"
)

type recorderRegistry struct {
	registry[recorder.IRecorder]
}

func (r *recorderRegistry) Register(name string, v recorder.IRecorder) error {
	return r.registry.Register(name, v)
}

func (r *recorderRegistry) Get(name string) recorder.IRecorder {
	if name != "" {
		return &recorderWrapper{name: name, r: r}
	}
	return nil
}

func (r *recorderRegistry) get(name string) recorder.IRecorder {
	return r.registry.Get(name)
}

type recorderWrapper struct {
	name string
	r    *recorderRegistry
}

func (w *recorderWrapper) Record(ctx context.Context, b []byte, opts ...recorder.RecordOption) error {
	v := w.r.get(w.name)
	if v == nil {
		return nil
	}
	return v.Record(ctx, b, opts...)
}
