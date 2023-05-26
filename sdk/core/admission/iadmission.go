package admission

import "context"

type IAdmission interface {
	Admit(ctx context.Context, addr string) bool
}

type admissionGroup struct {
	admissions []IAdmission
}

func AdmissionGroup(admissions ...IAdmission) IAdmission {
	return &admissionGroup{
		admissions: admissions,
	}
}

func (p *admissionGroup) Admit(ctx context.Context, addr string) bool {
	for _, admission := range p.admissions {
		if admission != nil && !admission.Admit(ctx, addr) {
			return false
		}
	}
	return true
}
