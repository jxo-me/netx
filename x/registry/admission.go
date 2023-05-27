package registry

import (
	"context"

	"github.com/jxo-me/netx/core/admission"
)

type admissionRegistry struct {
	registry[admission.IAdmission]
}

func (r *admissionRegistry) Register(name string, v admission.IAdmission) error {
	return r.registry.Register(name, v)
}

func (r *admissionRegistry) Get(name string) admission.IAdmission {
	if name != "" {
		return &admissionWrapper{name: name, r: r}
	}
	return nil
}

func (r *admissionRegistry) get(name string) admission.IAdmission {
	return r.registry.Get(name)
}

type admissionWrapper struct {
	name string
	r    *admissionRegistry
}

func (w *admissionWrapper) Admit(ctx context.Context, addr string) bool {
	p := w.r.get(w.name)
	if p == nil {
		return false
	}
	return p.Admit(ctx, addr)
}
