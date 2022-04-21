package wrapper

import (
	"bytes"
	"errors"
	"net"
	"syscall"

	"github.com/go-gost/core/limiter"
)

var (
	errUnsupport = errors.New("unsupported operation")
)

// serverConn is a server side Conn with metrics supported.
type serverConn struct {
	net.Conn
	rlimiter limiter.RateLimiter
	rbuf     bytes.Buffer
}

func WrapConn(rlimiter limiter.RateLimiter, c net.Conn) net.Conn {
	if rlimiter == nil {
		return c
	}
	return &serverConn{
		Conn:     c,
		rlimiter: rlimiter,
	}
}

func (c *serverConn) Read(b []byte) (n int, err error) {
	if c.rlimiter == nil || c.rlimiter.Input() == nil {
		return c.Conn.Read(b)
	}

	burst := len(b)
	if c.rbuf.Len() > 0 {
		if c.rbuf.Len() < burst {
			burst = c.rbuf.Len()
		}
		return c.rbuf.Read(b[:c.rlimiter.Input().Limit(burst)])
	}

	nn, err := c.Conn.Read(b)
	if err != nil {
		return nn, err
	}

	n = c.rlimiter.Input().Limit(nn)
	if n < nn {
		if _, err = c.rbuf.Write(b[n:nn]); err != nil {
			return 0, err
		}
	}

	return
}

func (c *serverConn) Write(b []byte) (n int, err error) {
	if c.rlimiter == nil || c.rlimiter.Output() == nil {
		return c.Conn.Write(b)
	}

	nn := 0
	for len(b) > 0 {
		nn, err = c.Conn.Write(b[:c.rlimiter.Output().Limit(len(b))])
		n += nn
		if err != nil {
			return
		}
		b = b[nn:]
	}

	return
}

func (c *serverConn) SyscallConn() (rc syscall.RawConn, err error) {
	if sc, ok := c.Conn.(syscall.Conn); ok {
		rc, err = sc.SyscallConn()
		return
	}
	err = errUnsupport
	return
}

type packetConn struct {
	net.PacketConn
	rlimiter limiter.RateLimiter
}

func WrapPacketConn(rlimiter limiter.RateLimiter, pc net.PacketConn) net.PacketConn {
	if rlimiter == nil {
		return pc
	}
	return &packetConn{
		PacketConn: pc,
		rlimiter:   rlimiter,
	}
}

func (c *packetConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	for {
		n, addr, err = c.PacketConn.ReadFrom(p)
		if err != nil {
			return
		}
		if c.rlimiter == nil || c.rlimiter.Input() == nil {
			return
		}

		if c.rlimiter.Input().Limit(n) < n {
			continue
		}

		return
	}
}

func (c *packetConn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	if c.rlimiter != nil &&
		c.rlimiter.Output() != nil &&
		c.rlimiter.Output().Limit(len(p)) < len(p) {
		n = len(p)
		return
	}

	return c.PacketConn.WriteTo(p, addr)
}
