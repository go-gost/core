package admission

import "context"

type Admission interface {
	Admit(ctx context.Context, addr string) bool
}

type admissionGroup struct {
	admissions []Admission
}

func AdmissionGroup(admissions ...Admission) Admission {
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
