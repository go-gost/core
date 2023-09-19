package auth

import "context"

// Authenticator is an interface for user authentication.
type Authenticator interface {
	Authenticate(ctx context.Context, user, password string) (ok bool, id string)
}

type authenticatorGroup struct {
	authers []Authenticator
}

func AuthenticatorGroup(authers ...Authenticator) Authenticator {
	return &authenticatorGroup{
		authers: authers,
	}
}

func (p *authenticatorGroup) Authenticate(ctx context.Context, user, password string) (bool, string) {
	if len(p.authers) == 0 {
		return false, ""
	}
	for _, auther := range p.authers {
		if auther == nil {
			continue
		}

		if ok, id := auther.Authenticate(ctx, user, password); ok {
			return ok, id
		}
	}
	return false, ""
}
