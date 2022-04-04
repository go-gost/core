package admission

type Admission interface {
	Admit(addr string) bool
}
