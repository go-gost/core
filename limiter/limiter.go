package limiter

// Limiter scopes
const (
	ScopeService string = "service"
	ScopeConn    string = "conn"
	ScopeClient  string = "client"
)

type Options struct {
	Service string
	Scope   string
	Network string
	Addr    string
	Client  string
	Src     string
}

type Option func(opts *Options)

func ServiceOption(service string) Option {
	return func(opts *Options) {
		opts.Service = service
	}
}

func ScopeOption(scope string) Option {
	return func(opts *Options) {
		opts.Scope = scope
	}
}

func NetworkOption(network string) Option {
	return func(opts *Options) {
		opts.Network = network
	}
}

func AddrOption(addr string) Option {
	return func(opts *Options) {
		opts.Addr = addr
	}
}

func ClientOption(client string) Option {
	return func(opts *Options) {
		opts.Client = client
	}
}

func SrcOption(src string) Option {
	return func(opts *Options) {
		opts.Src = src
	}
}
