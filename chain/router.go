package chain

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"time"

	"github.com/go-gost/core/hosts"
	"github.com/go-gost/core/logger"
	"github.com/go-gost/core/recorder"
	"github.com/go-gost/core/resolver"
)

type SockOpts struct {
	Mark int
}

type RouterOptions struct {
	Retries    int
	Timeout    time.Duration
	IfceName   string
	Netns      string
	SockOpts   *SockOpts
	Chain      Chainer
	Resolver   resolver.Resolver
	HostMapper hosts.HostMapper
	Recorders  []recorder.RecorderObject
	Logger     logger.Logger
}

type RouterOption func(*RouterOptions)

func InterfaceRouterOption(ifceName string) RouterOption {
	return func(o *RouterOptions) {
		o.IfceName = ifceName
	}
}

func NetnsRouterOption(netns string) RouterOption {
	return func(o *RouterOptions) {
		o.Netns = netns
	}
}

func SockOptsRouterOption(so *SockOpts) RouterOption {
	return func(o *RouterOptions) {
		o.SockOpts = so
	}
}

func TimeoutRouterOption(timeout time.Duration) RouterOption {
	return func(o *RouterOptions) {
		o.Timeout = timeout
	}
}

func RetriesRouterOption(retries int) RouterOption {
	return func(o *RouterOptions) {
		o.Retries = retries
	}
}

func ChainRouterOption(chain Chainer) RouterOption {
	return func(o *RouterOptions) {
		o.Chain = chain
	}
}

func ResolverRouterOption(resolver resolver.Resolver) RouterOption {
	return func(o *RouterOptions) {
		o.Resolver = resolver
	}
}

func HostMapperRouterOption(m hosts.HostMapper) RouterOption {
	return func(o *RouterOptions) {
		o.HostMapper = m
	}
}

func RecordersRouterOption(recorders ...recorder.RecorderObject) RouterOption {
	return func(o *RouterOptions) {
		o.Recorders = recorders
	}
}

func LoggerRouterOption(logger logger.Logger) RouterOption {
	return func(o *RouterOptions) {
		o.Logger = logger
	}
}

type Router struct {
	options RouterOptions
}

func NewRouter(opts ...RouterOption) *Router {
	r := &Router{}
	for _, opt := range opts {
		if opt != nil {
			opt(&r.options)
		}
	}
	if r.options.Logger == nil {
		r.options.Logger = logger.Default().WithFields(map[string]any{"kind": "router"})
	}
	return r
}

func (r *Router) Options() *RouterOptions {
	if r == nil {
		return nil
	}
	return &r.options
}

func (r *Router) Dial(ctx context.Context, network, address string) (conn net.Conn, err error) {
	host := address
	if h, _, _ := net.SplitHostPort(address); h != "" {
		host = h
	}
	r.record(ctx, recorder.RecorderServiceRouterDialAddress, []byte(host))

	conn, err = r.dial(ctx, network, address)
	if err != nil {
		r.record(ctx, recorder.RecorderServiceRouterDialAddressError, []byte(host))
		return
	}

	if network == "udp" || network == "udp4" || network == "udp6" {
		if _, ok := conn.(net.PacketConn); !ok {
			return &packetConn{conn}, nil
		}
	}
	return
}

func (r *Router) record(ctx context.Context, name string, data []byte) error {
	if len(data) == 0 {
		return nil
	}

	for _, rec := range r.options.Recorders {
		if rec.Record == name {
			err := rec.Recorder.Record(ctx, data)
			if err != nil {
				r.options.Logger.Errorf("record %s: %v", name, err)
			}
			return err
		}
	}
	return nil
}

func (r *Router) dial(ctx context.Context, network, address string) (conn net.Conn, err error) {
	count := r.options.Retries + 1
	if count <= 0 {
		count = 1
	}
	r.options.Logger.Debugf("dial %s/%s", address, network)

	for i := 0; i < count; i++ {
		var ipAddr string
		ipAddr, err = Resolve(ctx, "ip", address, r.options.Resolver, r.options.HostMapper, r.options.Logger)
		if err != nil {
			r.options.Logger.Error(err)
			break
		}

		var route Route
		if r.options.Chain != nil {
			route = r.options.Chain.Route(ctx, network, ipAddr, WithHostRouteOption(address))
		}

		if r.options.Logger.IsLevelEnabled(logger.DebugLevel) {
			buf := bytes.Buffer{}
			for _, node := range routePath(route) {
				fmt.Fprintf(&buf, "%s@%s > ", node.Name, node.Addr)
			}
			fmt.Fprintf(&buf, "%s", ipAddr)
			r.options.Logger.Debugf("route(retry=%d) %s", i, buf.String())
		}

		if route == nil {
			route = DefaultRoute
		}
		conn, err = route.Dial(ctx, network, ipAddr,
			InterfaceDialOption(r.options.IfceName),
			NetnsDialOption(r.options.Netns),
			SockOptsDialOption(r.options.SockOpts),
			LoggerDialOption(r.options.Logger),
			TimeoutDialOption(r.options.Timeout),
		)
		if err == nil {
			break
		}
		r.options.Logger.Errorf("route(retry=%d) %s", i, err)
	}

	return
}

func (r *Router) Bind(ctx context.Context, network, address string, opts ...BindOption) (ln net.Listener, err error) {
	count := r.options.Retries + 1
	if count <= 0 {
		count = 1
	}
	r.options.Logger.Debugf("bind on %s/%s", address, network)

	for i := 0; i < count; i++ {
		var route Route
		if r.options.Chain != nil {
			route = r.options.Chain.Route(ctx, network, address)
			if route == nil || len(route.Nodes()) == 0 {
				err = ErrEmptyRoute
				return
			}
		}

		if r.options.Logger.IsLevelEnabled(logger.DebugLevel) {
			buf := bytes.Buffer{}
			for _, node := range routePath(route) {
				fmt.Fprintf(&buf, "%s@%s > ", node.Name, node.Addr)
			}
			fmt.Fprintf(&buf, "%s", address)
			r.options.Logger.Debugf("route(retry=%d) %s", i, buf.String())
		}

		if route == nil {
			route = DefaultRoute
		}
		ln, err = route.Bind(ctx, network, address, opts...)
		if err == nil {
			break
		}
		r.options.Logger.Errorf("route(retry=%d) %s", i, err)
	}

	return
}

func routePath(route Route) (path []*Node) {
	if route == nil {
		return
	}
	for _, node := range route.Nodes() {
		if tr := node.Options().Transport; tr != nil {
			path = append(path, routePath(tr.Options().Route)...)
		}
		path = append(path, node)
	}
	return
}

type packetConn struct {
	net.Conn
}

func (c *packetConn) ReadFrom(b []byte) (n int, addr net.Addr, err error) {
	n, err = c.Read(b)
	addr = c.Conn.RemoteAddr()
	return
}

func (c *packetConn) WriteTo(b []byte, addr net.Addr) (n int, err error) {
	return c.Write(b)
}
