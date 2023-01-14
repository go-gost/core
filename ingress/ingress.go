package ingress

type Ingress interface {
	Get(host string) string
}
