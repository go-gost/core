package auth

import "context"

// Authenticator is an interface for user authentication.
type Authenticator interface {
	Authenticate(ctx context.Context, user, password string) (id string, ok bool)
}

type authenticatorGroup struct {
	authers []Authenticator
}

func AuthenticatorGroup(authers ...Authenticator) Authenticator {
	return &authenticatorGroup{
		authers: authers,
	}
}

func (p *authenticatorGroup) Authenticate(ctx context.Context, user, password string) (string, bool) {
	if len(p.authers) == 0 {
		return "", false
	}
	for _, auther := range p.authers {
		if auther == nil {
			continue
		}

		if id, ok := auther.Authenticate(ctx, user, password); ok {
			return id, ok
		}
	}
	return "", false
}
