package auth

// Authenticator is an interface for user authentication.
type Authenticator interface {
	Authenticate(user, password string) bool
}

type authenticatorGroup struct {
	authers []Authenticator
}

func AuthenticatorGroup(authers ...Authenticator) Authenticator {
	return &authenticatorGroup{
		authers: authers,
	}
}

func (p *authenticatorGroup) Authenticate(user, password string) bool {
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
