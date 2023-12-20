package admission

import "context"

type Options struct{}
type Option func(opts *Options)

type IAdmission interface {
	Admit(ctx context.Context, addr string, opts ...Option) bool
}

type admissionGroup struct {
	admissions []IAdmission
}

func AdmissionGroup(admissions ...IAdmission) IAdmission {
	return &admissionGroup{
		admissions: admissions,
	}
}

func (p *admissionGroup) Admit(ctx context.Context, addr string, opts ...Option) bool {
	for _, admission := range p.admissions {
		if admission != nil && !admission.Admit(ctx, addr, opts...) {
			return false
		}
	}
	return true
}
