package chain

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"time"

	"github.com/go-gost/core/connector"
	"github.com/go-gost/core/hosts"
	"github.com/go-gost/core/logger"
	"github.com/go-gost/core/recorder"
	"github.com/go-gost/core/resolver"
)

type SockOpts struct {
	Mark int
}

type Router struct {
	ifceName  string
	sockOpts  *SockOpts
	timeout   time.Duration
	retries   int
	chain     Chainer
	resolver  resolver.Resolver
	hosts     hosts.HostMapper
	recorders []recorder.RecorderObject
	logger    logger.Logger
}

func (r *Router) WithTimeout(timeout time.Duration) *Router {
	r.timeout = timeout
	return r
}

func (r *Router) WithRetries(retries int) *Router {
	r.retries = retries
	return r
}

func (r *Router) WithInterface(ifceName string) *Router {
	r.ifceName = ifceName
	return r
}

func (r *Router) WithSockOpts(so *SockOpts) *Router {
	r.sockOpts = so
	return r
}

func (r *Router) WithChain(chain Chainer) *Router {
	r.chain = chain
	return r
}

func (r *Router) WithResolver(resolver resolver.Resolver) *Router {
	r.resolver = resolver
	return r
}

func (r *Router) WithHosts(hosts hosts.HostMapper) *Router {
	r.hosts = hosts
	return r
}

func (r *Router) Hosts() hosts.HostMapper {
	if r != nil {
		return r.hosts
	}
	return nil
}

func (r *Router) WithRecorder(recorders ...recorder.RecorderObject) *Router {
	r.recorders = recorders
	return r
}

func (r *Router) WithLogger(logger logger.Logger) *Router {
	r.logger = logger
	return r
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

	for _, rec := range r.recorders {
		if rec.Record == name {
			err := rec.Recorder.Record(ctx, data)
			if err != nil {
				r.logger.Errorf("record %s: %v", name, err)
			}
			return err
		}
	}
	return nil
}

func (r *Router) dial(ctx context.Context, network, address string) (conn net.Conn, err error) {
	count := r.retries + 1
	if count <= 0 {
		count = 1
	}
	r.logger.Debugf("dial %s/%s", address, network)

	for i := 0; i < count; i++ {
		var route *Route
		if r.chain != nil {
			route = r.chain.Route(network, address)
		}

		if r.logger.IsLevelEnabled(logger.DebugLevel) {
			buf := bytes.Buffer{}
			for _, node := range route.Path() {
				fmt.Fprintf(&buf, "%s@%s > ", node.Name, node.Addr)
			}
			fmt.Fprintf(&buf, "%s", address)
			r.logger.Debugf("route(retry=%d) %s", i, buf.String())
		}

		address, err = resolve(ctx, "ip", address, r.resolver, r.hosts, r.logger)
		if err != nil {
			r.logger.Error(err)
			break
		}

		if route == nil {
			route = &Route{}
		}
		route.logger = r.logger
		conn, err = route.Dial(ctx, network, address,
			InterfaceDialOption(r.ifceName),
			SockOptsDialOption(r.sockOpts),
		)
		if err == nil {
			break
		}
		r.logger.Errorf("route(retry=%d) %s", i, err)
	}

	return
}

func (r *Router) Bind(ctx context.Context, network, address string, opts ...connector.BindOption) (ln net.Listener, err error) {
	count := r.retries + 1
	if count <= 0 {
		count = 1
	}
	r.logger.Debugf("bind on %s/%s", address, network)

	for i := 0; i < count; i++ {
		var route *Route
		if r.chain != nil {
			route = r.chain.Route(network, address)
			if route.Len() == 0 {
				err = ErrEmptyRoute
				return
			}
		}

		if r.logger.IsLevelEnabled(logger.DebugLevel) {
			buf := bytes.Buffer{}
			for _, node := range route.Path() {
				fmt.Fprintf(&buf, "%s@%s > ", node.Name, node.Addr)
			}
			fmt.Fprintf(&buf, "%s", address)
			r.logger.Debugf("route(retry=%d) %s", i, buf.String())
		}

		ln, err = route.Bind(ctx, network, address, opts...)
		if err == nil {
			break
		}
		r.logger.Errorf("route(retry=%d) %s", i, err)
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
