package auth

import "context"

// Authenticator is an interface for user authentication.
type Authenticator interface {
	Authenticate(ctx context.Context, user, password string) bool
}

type authenticatorGroup struct {
	authers []Authenticator
}

func AuthenticatorGroup(authers ...Authenticator) Authenticator {
	return &authenticatorGroup{
		authers: authers,
	}
}

func (p *authenticatorGroup) Authenticate(ctx context.Context, user, password string) bool {
	if len(p.authers) == 0 {
		return true
	}
	for _, auther := range p.authers {
		if auther != nil && auther.Authenticate(ctx, user, password) {
			return true
		}
	}
	return false
}
