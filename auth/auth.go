package auth

// Authenticator is an interface for user authentication.
type Authenticator interface {
	Authenticate(user, password string) bool
}

type authenticatorList struct {
	authers []Authenticator
}

func AuthenticatorList(authers ...Authenticator) Authenticator {
	return &authenticatorList{
		authers: authers,
	}
}

func (p *authenticatorList) Authenticate(user, password string) bool {
	if len(p.authers) == 0 {
		return true
	}
	for _, auther := range p.authers {
		if auther != nil && auther.Authenticate(user, password) {
			return true
		}
	}
	return false
}
