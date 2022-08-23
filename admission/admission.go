package admission

type Admission interface {
	Admit(addr string) bool
}

type admissionList struct {
	admissions []Admission
}

func AdmissionList(admissions ...Admission) Admission {
	return &admissionList{
		admissions: admissions,
	}
}

func (p *admissionList) Admit(addr string) bool {
	for _, admission := range p.admissions {
		if admission != nil && !admission.Admit(addr) {
			return false
		}
	}
	return true
}
