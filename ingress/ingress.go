package ingress

import "context"

type Ingress interface {
	Get(ctx context.Context, host string) string
}
