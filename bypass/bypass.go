package bypass

// Bypass is a filter of address (IP or domain).
type Bypass interface {
	// Contains reports whether the bypass includes addr.
	Contains(addr string) bool
}

type bypassList struct {
	bypasses []Bypass
}

func BypassList(bypasses ...Bypass) Bypass {
	return &bypassList{
		bypasses: bypasses,
	}
}

func (p *bypassList) Contains(addr string) bool {
	for _, bypass := range p.bypasses {
		if bypass != nil && bypass.Contains(addr) {
			return true
		}
	}
	return false
}
