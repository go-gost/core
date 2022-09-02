package admission

type Admission interface {
	Admit(addr string) bool
}

type admissionGroup struct {
	admissions []Admission
}

func AdmissionGroup(admissions ...Admission) Admission {
	return &admissionGroup{
		admissions: admissions,
	}
}

func (p *admissionGroup) Admit(addr string) bool {
	for _, admission := range p.admissions {
		if admission != nil && !admission.Admit(addr) {
			return false
		}
	}
	return true
}
