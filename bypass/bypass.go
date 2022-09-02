package bypass

// Bypass is a filter of address (IP or domain).
type Bypass interface {
	// Contains reports whether the bypass includes addr.
	Contains(addr string) bool
}

type bypassGroup struct {
	bypasses []Bypass
}

func BypassGroup(bypasses ...Bypass) Bypass {
	return &bypassGroup{
		bypasses: bypasses,
	}
}

func (p *bypassGroup) Contains(addr string) bool {
	for _, bypass := range p.bypasses {
		if bypass != nil && bypass.Contains(addr) {
			return true
		}
	}
	return false
}
